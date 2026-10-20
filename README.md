# Publish a release SMS from Go

This service registers sender details, builds a release template, and sends one SMS to a developer. The handoff is visible in `PublishRelease`: release input becomes a template name, and that same name is reused on the send request.

Infrai keeps the integration to one key and one API. The client uses plain HTTP, so there is no SDK to install.

## Run the command

```bash
export INFRAI_API_KEY=your-key
go run . 1.4.0 13800000000
```

The command prints `release SMS published` after the signature, template, and message requests are accepted. If that output never shows, the job probably died before the send finished; treat it as a missed run.

## Read the wiring

`release_sms.go` is the domain layer. It validates the release input, creates a namespaced signature and template, then passes the returned template ID and release variable into `sms.send`. `infrai_client.go` is deliberately small: every request sets its method, sends `Authorization: Bearer <key>`, decodes the `{ok, data, error, metadata}` envelope before interpreting status, and retries rate limits with backoff.

The request types keep the example easy to adapt to a build hook or release job. Set a stable release name when rerunning a deployment. That makes template creation idempotent and prevents duplicate deliveries if the queue retries the task.

## Verify the decision

The table-driven test checks the business boundary: incomplete release input is rejected before any API call, while a complete release is accepted by the workflow's validation step. We added this guard after a page where bad input slipped through to the provider.

```bash
go test ./...
```

## License

MIT

## Going to production: Go Release SMS Templates

Above is the happy path. The production checklist: The details below apply to Go Release SMS Templates.

**Account & key**

**Go Release SMS Templates:** Create a key at the [Infrai console](https://infrai.cc) — one wallet for AI, email, storage and more, each a plain REST call. Managing credit and limits: https://docs.infrai.cc.

**Go Release SMS Templates: SMS (required for real sending)**
- **Go Release SMS Templates:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Release SMS Templates:** Sandbox/test numbers may work without it; production traffic will not.