# Sensuctl Plugin Template

## Overview
sensuctl-plugin-template is a template repository which wraps the [Sensu Plugin SDK][2].
To use this project as a template, click the "Use this template" button from the main project page.
Once the repository is created from this template, you can use the [Sensu Plugin Tool][9] to
populate the templated fields with the proper values.

## Functionality

After successfully creating a project from this template, update [main.go][7] to customize the
behavior of this sensuctl plugin.

When writing or updating a plugin's README from this template, review the Sensu Community
[plugin README style guide][3] for content suggestions and guidance. Remove everything
prior to `# namespace-copy` from the generated README file, and add additional context about the
plugin per the style guide.

## Releases with Github Actions

To release a version of your project, simply tag the target sha with a semver release without a `v`
prefix (ex. `1.0.0`). This will trigger the [GitHub action][5] workflow to [build and release][4]
the plugin with goreleaser. Register the asset with [Bonsai][8] to share it with the community!

***

# namespace-copy

## Table of Contents
- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Resource definition](#resource-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

The namespace-copy is a [Sensu CLI][6] tool (`sensuctl`) that ...

## Files

## Usage examples

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add jspaleta/sensu-namespace-copy
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/jspaleta/sensu-namespace-copy].

### Resource definition

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-namespace-copy repository:

```
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu-community/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sensu-community/sensuctl-plugin-template/blob/master/.github/workflows/release.yml
[5]: https://github.com/sensu-community/sensuctl-plugin-template/actions
[6]: https://docs.sensu.io/sensu-go/latest/sensuctl/reference/
[7]: https://github.com/sensu-community/sensuctl-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu-community/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
