open Lwt.Infix

exception Closed

let rec read_exactly fd p len =
  (* TODO: socket conditions e.g. EAGAIN? *)
  Unistd_unix_lwt.read fd (Ctypes.to_voidp p) len
  >>= function
  | 0 -> Lwt.fail Closed
  | n ->
    let remaining = len - n in
    if remaining = 0
    then Lwt.return_unit
    else read_exactly fd Ctypes.(p +@ n) remaining

let read_message fd =
  let len_ptr = Ctypes.(allocate uint32_t) Unsigned.UInt32.zero in
  read_exactly fd len_ptr 4
  >>= fun () ->
  let len = Ctypes.(!@ len_ptr) in
  let count = Unsigned.UInt32.to_int len in
  let buf = Ctypes.(allocate_n uint8_t) ~count in
  Ctypes.((coerce (ptr uint8_t) (ptr uint32_t) buf) <-@ len);
  read_exactly fd Ctypes.(buf +@ 4) (count - 4)
  >>= fun () ->
  Lwt.return (Ctypes.CArray.from_ptr buf (Unsigned.UInt32.to_int len))

let rec write_exactly fd p len =
  (* TODO: socket conditions e.g. EAGAIN? *)
  Unistd_unix_lwt.write fd (Ctypes.to_voidp p) len
  >>= function
  | 0 -> Lwt.fail Closed
  | n ->
    let remaining = len - n in
    if remaining = 0
    then Lwt.return_unit
    else write_exactly fd Ctypes.(p +@ n) remaining
