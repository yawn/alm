# AWS last messages (`alm`)

`alm` continiously fetches the last (as identified by timestamp and a configurable delta from _now_) messages from various services and persists them in a directory as logfiles. It currently only supports the AWS standard partition and queries all active regions for the given services. Multi account discovery is currently not available but planned (probably pending on AWS SSO support in the v2 SDK).

The following services are supported (PRs welcome):

- CloudFormation (Stacks and StackSets): this generates artificial logfiles over a selected subset of stack events using the convention `cfr-:account-:region:-:stackid.log`
- CloudTrail: fetches individual log groups using the convention `cwl-:account-:region:-:logstream-id-:loggroup-id.log`

When files exist, `alm` will delete them.
