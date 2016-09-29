module ForMac = struct
  type t  =
    | Stable1_12_0
    | Stable1_12_0_a
    | Stable1_12_1
    | Beta21
    | Beta22
    | Beta23
    | Beta24
    | Beta25
    | Beta26
    | Beta26_1

  let to_string = function
    | Stable1_12_0 -> "version/1.12.0"
    | Stable1_12_0_a -> "version/1.12.0-a"
    | Stable1_12_1 -> "version/1.12.1"
    | Beta21 -> "version/beta21"
    | Beta22 -> "version/beta22"
    | Beta23 -> "version/beta23"
    | Beta24 -> "version/beta24"
    | Beta25 -> "version/beta25"
    | Beta26 -> "version/beta26"
    | Beta26_1 -> "version/beta26.1"

  let of_string = function
    | "1.12.0" -> Some Stable1_12_0
    | "1.12.0-a" -> Some Stable1_12_0_a
    | "1.12.1" -> Some Stable1_12_1
    | "1.12.0-beta21" -> Some Beta21
    | "1.12.0-beta22" -> Some Beta22
    | "1.12.1-rc1-beta23" -> Some Beta23
    | "1.12.1-beta24" -> Some Beta24
    | "1.12.1-beta25" -> Some Beta25
    | "1.12.1-beta26" -> Some Beta26
    | "1.12.1-beta26.1" -> Some Beta26_1
    | _ -> None
end

module Osx = struct
  type t =
    | Yosemite
    | ElCapitan
    | Sierra

  let to_string = function
    | Yosemite -> "osx/10.10.x"
    | ElCapitan -> "osx/10.11.x"
    | Sierra -> "osx/10.12 (beta)"

  let of_string = function
    | "10.10" -> Some Yosemite
    | "10.11" -> Some ElCapitan
    | "10.12" -> Some Sierra
    | _ -> None
end

module Status = struct
  (* type t should be in Stage order (see Stage module below) *)
  type t =
    | NeedsTriage
    | MoreInfoNeeded
    | WontFix
    | Acknowledged
    | InProgress
    | Fixed
    | ReleasedBeta
    | ReleasedStable

  module Set = Set.Make(struct
      type nonrec t = t

      let compare = compare
    end)

  module Stage = struct
    type t =
      | TriageStage
      | WontFixStage
      | AcknowledgedStage
      | InProgressStage
      | FixedStage
      | ReleasedStage

    let of_status = function
      | NeedsTriage | MoreInfoNeeded -> TriageStage
      | WontFix -> WontFixStage
      | Acknowledged -> AcknowledgedStage
      | InProgress -> InProgressStage
      | Fixed -> FixedStage
      | ReleasedBeta | ReleasedStable -> ReleasedStage

    let to_string = function
      | TriageStage -> "triage"
      | WontFixStage -> "wontfix"
      | AcknowledgedStage -> "acknowledged"
      | InProgressStage -> "inprogress"
      | FixedStage -> "fixed"
      | ReleasedStage -> "released"
  end

  let to_string = function
    | NeedsTriage -> "status/0-triage"
    | WontFix -> "status/0-wont-fix"
    | MoreInfoNeeded -> "status/0-more-info-needed"
    | Acknowledged -> "status/1-acknowledged"
    | InProgress -> "status/2-in-progress"
    | Fixed -> "status/3-fixed"
    | ReleasedBeta -> "status/4-fix-released-beta"
    | ReleasedStable -> "status/4-fix-released-stable"

  let of_string = function
    | "status/0-triage" -> Some NeedsTriage
    | "status/0-wont-fix" -> Some WontFix
    | "status/0-more-info-needed" -> Some MoreInfoNeeded
    | "status/1-acknowledged" -> Some Acknowledged
    | "status/2-in-progress" -> Some InProgress
    | "status/3-fixed" -> Some Fixed
    | "status/4-fix-released-beta" -> Some ReleasedBeta
    | "status/4-fix-released-stable" -> Some ReleasedStable
    | _ -> None
end

type t = { for_mac : ForMac.t option;  osx : Osx.t option; status : Status.t }

let to_string_list l =
  let result = [Status.to_string l.status] in
  let result = match l.for_mac with
  | None -> result
  | Some x -> result@ [ForMac.to_string x]
  in
  match l.osx with
  | None -> result
  | Some x -> result@ [Osx.to_string x]

let to_string l =
  let list_of_l = to_string_list l in
  String.concat ", " list_of_l

