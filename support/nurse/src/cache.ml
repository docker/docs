open Lwt.Infix

let pinata_fetch = "pinata:fetch" (* image used to fetch from S3 *)

type 'a error = ('a, [ `Msg of string ]) Result.result

type diagnostic_id = string

type timestamp = string

let (/) = Filename.concat

let home =
  try
    Unix.getenv "HOME"
  with Not_found ->
    Logs.err (fun f -> f "Please set the HOME environment variable");
    exit 1

let default_cache_dir =
  let dir = home / ".nurse" in
  (try Unix.mkdir dir 0o0700 with Unix.Unix_error(Unix.EEXIST, _, _) -> ());
  dir

let errorf fmt =
  Printf.ksprintf (fun s ->
    Result.Error (`Msg s)
  ) fmt

let path_of diagnostic_id timestamp =
  default_cache_dir / diagnostic_id / timestamp / "decompressed.tar"

let fetch diagnostic_id =
  (* Check AWS credentials have been setup *)
  let credentials_path = home / ".aws" / "credentials" in
  if not(Sys.file_exists credentials_path)
  then Lwt.return (errorf "Failed to find AWS credentials in %s. Please follow the instructions in the README.md" credentials_path)
  else begin
    Docker.inspect pinata_fetch
    >>= function
    | Result.Error (`Msg m) -> Lwt.return (Result.Error (`Msg m))
    | Result.Error (`Exit _n) -> Lwt.return (errorf "Failed to find image %s. Please build it using the instructions in the README.md" pinata_fetch)
    | Result.Ok _ ->
      let log_volume = default_cache_dir ^ ":/logs" in
      let cred_volume = credentials_path ^ ":/root/.aws/credentials" in
      Docker.run [ "-it"; "--rm"; "-v"; log_volume; "-v"; cred_volume; pinata_fetch; "mac"; diagnostic_id ]
      >>= function
      | Result.Error (`Msg m) -> Lwt.return (Result.Error (`Msg m))
      | Result.Error (`Exit _n) -> Lwt.return (errorf "Failed to fetch diagnostic ID %s. Please check for typos. Perhaps the upload never completed?" diagnostic_id)
      | Result.Ok _ ->
        (* Find and decompress *)
        let timestamps = Array.to_list @@ Sys.readdir (default_cache_dir / diagnostic_id) in
        Lwt_list.iter_s
          (fun timestamp ->
            let dir = default_cache_dir / diagnostic_id / timestamp in
            if Sys.is_directory dir then begin
              (* Freshly downloaded files are compressed, so fix the name *)
              ( if Sys.file_exists (dir / "report.tar") then begin
                  Lwt_unix.rename (dir / "report.tar") (dir / "decompressed.tar.gz")
                end else Lwt.return () )
              >>= fun () ->
              (* Properly-named compressed files should be decompressed *)
              ( if Sys.file_exists (dir / "decompressed.tar.gz") then begin
                  Command.run "/usr/bin/gunzip" [ dir / "decompressed.tar.gz" ]
                  >>= function
                  | Result.Error (`Exit 1) ->
                    (* For some reason the archives have trailing garbage. Is
                       this caused by an S3 upload block size? We consider it a
                       success if the decompressed file is created. *)
                    if Sys.file_exists (dir / "decompressed.tar")
                    then Lwt.return ()
                    else begin
                      Logs.err (fun f -> f "Failed to gunzip %s" (dir / "decompressed.tar.gz"));
                      Lwt.return ()
                    end
                  | Result.Error (`Exit _n) ->
                    Logs.err (fun f -> f "Failed to gunzip %s" (dir / "decompressed.tar.gz"));
                    Lwt.return ()
                  | Result.Error (`Msg m) ->
                    Logs.err (fun f -> f "%s" m);
                    Lwt.return ()
                  | Result.Ok _ ->
                    Lwt.return ()
                end else Lwt.return () )
            end else Lwt.return ()
          ) timestamps
        >>= fun () ->
        Lwt.return (Result.Ok ())
  end

type path = string

let get diagnostic_id timestamp =
  let path = path_of diagnostic_id timestamp in
  if Sys.file_exists path
  then Lwt.return (Result.Ok path)
  else begin
    Logs.info (fun f -> f "Archive not found in cache, asking S3");
    fetch diagnostic_id
    >>= function
    | Result.Error (`Msg m) ->
      Lwt.return (Result.Error (`Msg m))
    | Result.Ok () ->
      if Sys.file_exists path
      then Lwt.return (Result.Ok path)
      else Lwt.return (errorf "Failed to find archive %s/%s on S3" diagnostic_id timestamp)
  end
