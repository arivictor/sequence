![](./cover.png)

Sequence is a flexible job execution tool designed to manage and run a series of tasks defined in a YAML configuration file. Sequence together conditional execution of tasks defined in-line, bash scripts, python, javascript, or anything available in your terminal!

```shell
» sequence --config config.yml

2024/01/19 12:49:28 execute : 'Job 1'
Hello Exit Status 1!

2024/01/19 12:49:30 execute : 'Job 1' : error_handler 'handle_job_1_error'
Job failed, handling error..

2024/01/19 12:49:30 error : 'Job 1' : exit status 1
2024/01/19 12:49:30 skip : 'Job 2' : skip is true
2024/01/19 12:49:30 pass : 'Job 3' : depends_on 'Job 1'
2024/01/19 12:49:30 execute : 'Job 4'
This will run..

2024/01/19 12:49:30 execute : 'Job 4' : success_handler 'handle_job_4_success'
Job 4 succeeded, handling success..
```

## Features

- **Job Execution**: Execute a series of jobs defined in a structured YAML file.
- **Error Handling**: Specify error handlers for each job, providing robust control over error management.
- **Conditional Execution**: Define dependencies among jobs, ensuring that certain jobs run only after their dependencies have executed successfully.
- **Job Skipping**: Flexibly skip certain jobs without removing them from the configuration, allowing for dynamic adjustments of the job sequence.
- **Exit Control**: Control the continuation or termination of the sequence when a job fails, based on the `exit_on_error` attribute.

## Getting Started

### Prerequisites

- Go (version 1.15 or later)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/arivictor/sequence.git
    ```

2. Navigate to the repository directory:

    ```bash
    cd sequence
    ```

3. Build the application:

    ```bash
    go build -o sequence ./cmd/sequence
    ```

### Configuration Schema

Define your jobs and their properties in a YAML file. Below is a template of how the configuration should look:

```yaml
jobs:
  - name: "Job 1"
    command: "echo 'Hello World!' && exit 1"
    exit_on_error: false # quits the process on error
    on_error: "handle_err"
    on_success: "handle_success"
    skip: false

  - name: "Job 2"
    command: "node hello_world.js"
    depends_on: ["Job 1"] # Won't run, depends on Job 1

success_handlers:
    - name: "handle_success"
      command: "echo 'handling success..'"

error_handlers:
  - name: "handle_err"
    command: "echo 'handling an error..'"
```

> [!NOTE]  
> See [example.yaml](./example.yaml) for a more verbose example..


### Running Sequence

To execute the jobs as per your configuration file (e.g., `config.yaml`), use the following command:

```bash
./sequence -config ./config.yaml
```

> [!NOTE]  
> Filepaths defined in jobs are executed relative from where the command is executed.

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.