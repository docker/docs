type change =
  | Close of string
  | Remove_labels of string * string list

type action = {
  issue : Github_t.issue;
  changes : change list;
}

let string_of_change = function
  | Close descr -> "Close because "^descr
  | Remove_labels (descr, labels) ->
    Printf.sprintf "Remove labels [%s] because %s"
      (String.concat ", " labels) descr

(* TODO: Offer to close issues (with a note) that are very old or
   otherwise outdated. *)
let compute_actions list issue =
  let open Github_t in
  let open Labels in

  let statuses = List.fold_left (fun set label ->
    match Status.of_string label.label_name with
    | Some Status.MoreInfoNeeded (* MoreInfoNeeded is always allowed *)
    | None -> set
    | Some status -> Status.Set.add status set
  ) Status.Set.empty issue.issue_labels in

  (* Compute status labels to remove due to later stage *)
  let latest_status =
    try Status.Set.max_elt statuses with Not_found -> Status.NeedsTriage
  in
  let latest_stage = Status.Stage.of_status latest_status in
  let _to_keep, to_remove = Status.Set.partition (fun status ->
    latest_stage = Status.Stage.of_status status
  ) statuses in
  let changes =
    if Status.Set.is_empty to_remove then []
    else
      let message = Printf.sprintf "issue is at stage %s"
          (Status.Stage.to_string latest_stage)
      in
      let labels =
        List.map Status.to_string (Status.Set.elements to_remove)
      in
      [ Remove_labels (message, labels) ]
  in

  (* Decide whether to close issue due to completion *)
  let changes =
    if Status.(Set.mem ReleasedBeta) statuses
    && Status.(Set.mem ReleasedStable) statuses
    then
      let message = "a fix has been released on beta and stable channels" in
      (Close message)::changes
    else changes
  in

  (* Decide whether to close issue due to WontFix *)
  let changes =
    if Status.(Set.mem WontFix) statuses
    then (Close "the issue has been marked WontFix")::changes
    else changes
  in

  Github.Monad.return (match changes with
    | [] -> list
    | changes -> { issue; changes; }::list
  )

let clean_issue ~user ~repo action =
  let open Github.Monad in
  let num = action.issue.Github_t.issue_number in
  List.fold_left (fun m -> function
    | Close _ ->
      m >>= fun () ->
      let issue = Github_t.{
        update_issue_title = None;
        update_issue_body  = None;
        update_issue_state = Some `Closed;
        update_issue_assignee = None;
        update_issue_milestone = None;
        update_issue_labels = None;
      } in
      Github.Issue.update ~user ~repo ~num ~issue ()
      >>~ fun _updated_issue ->
      return ()
    | Remove_labels (_, labels) -> List.fold_left (fun m name ->
      m >>= fun () ->
      Github.Issue.remove_label ~user ~repo ~num ~name ()
      >>~ fun _labels ->
      return ()
    ) m labels
  ) (return ()) action.changes
