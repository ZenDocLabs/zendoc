# Zendoc CLI Documentation

Zendoc CLI offers two main commands that you can use to interact with the tool.

## Init Command

```bash
zendoc init
```

This command initializes a Zendoc project within your Go project. It creates a `.zendoc.config.json` file in the root of your project. Here is an example of this file:

```json
{
  "projectConfig": {
    "name": "zendoc",
    "description": "Description of your project",
    "version": "1.0",
    "gitLink": "https://github.com/ZenDocLabs/zendoc",
    "mainBranch": "main",
    "docPath": "./doc"
  },
  "docConfig": {
    "includePrivate": false,
    "includeTests": false,
    "includeMain": false,
    "excludeFiles": []
  }
}
```

This JSON file contains two main configuration sections:

### 1. projectConfig

This section contains information about your project:

- `name`: the name of your project
- `description`: a brief description of what your project does
- `version`: the current version of your documentation (critical value for documentation versioning)
- `gitLink`: the Git repository link of your project
- `mainBranch`: the main branch used in your project
- `docPath`: the path where the documentation web application will be created

### 2. docConfig

This section configures the CLI behavior, particularly for the `generate` command:

- `includePrivate`: enables documentation generation for private functions
- `includeTests`: enables documentation generation for test files
- `includeMain`: feature not currently in use (likely to be implemented soon)
- `excludeFiles`: array of regular expressions to exclude matching files during documentation generation

## Generate Command

```bash
zendoc generate <output>
```

The `output` parameter can take two values: `json` or `web`.

### `json` Option

The command analyzes your documentation and exports it to a file named `doc.json`.

### `web` Option

The command performs the following operations:

1. Clones the [template](https://github.com/ZenDocLabs/zendoc-ui-template) to create a web interface
2. Analyzes your documentation and exports it to a file named `data-[your-version].json` in the `doc/src/assets` folder
3. Creates an `.env` file containing your application name and Git information
4. Updates the `versions.json` file located next to the data files

Documentation versioning is managed through the `version` value in your `.zendoc.config.json` file. To create multiple documentation versions, simply change this value.
Once the `web` option is used, you can simply go to the generated web-app and run `npm run dev` to see the beautiful result !
