# apigen

This repository contains a sample project how to utilize go:generate in order to overcome the burden of writing boilerplate code. It contains a REST API and a couple of middlewares. The tool utilizing go:generate automatically creates code for the following functionality:

- Logging
- Instrumentation
- Activity Recording
- API Documentation

This code was partially shown and demonstrated in the GopherCon Russia 2019 talk **go generate: One File To Rule Them All**.

[Slides](https://speakerdeck.com/konradreiche/go-generate-one-file-to-rule-them-all)
[Video](https://www.youtube.com/watch?v=RfKgBI4JgSI)

You can run the code generation with:

```bash
go generate ./...
```
