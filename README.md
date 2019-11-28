# KPaaS

KPaaS is a tool for deploying Kubernetes clusters, operating Kubernetes clusters, and managing software lifecycles in Kubernetes clusters.

## Quick Start

```shell script
# Clone code
git clone https://github.com/kpaas-io/kpaas.git

# Build and run
make run
```

### Go to [http://localhost:8080](http://localhost:8080)

Let's deploy a Kubernetes Cluster.

## Features

* deploy a Kubernetes Cluster

## Documentation

Not Yet

## Development

### Auto Add License Headers

#### VSCode

Install the [licenser](https://marketplace.visualstudio.com/items?itemName=ymotongpoo.licenser) extension. It will insert license header automatically when creatiing a new file. Or you can manually add license header using "licenser: Insert license header" commands via Command Palette (`Ctrl+Shift+P` on Windows/Linux, `⌘⇧P` on OS X)

Add editor's settings in `kpaas/.vscode/settings.json`:

```json
{
  ...,
  "licenser.projectName": "KPaaS",
  "licenser.author": "Shanghai JingDuo Information Technology co., Ltd.",
  ...
}
```

#### JetBrain

1. In the **Settings/Preferences** dialog, select **Editor** -> **Copyright** -> **Copyright Profiles**.
2. Click Add and name the new profile.
3. Enter copyright notice text.
```text
Copyright $today.year Shanghai JingDuo Information Technology co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
4. In the **Settings/Preferences** dialog, go to **Appearance & Behavior** | **Scopes**.
5. Click **+** Add scope (Insert) icon, select **Shared** from the list, and specify the name of the scope.
6. Specify a pattern in the Pattern field manually, fill `*&&!file[kpaas]:*&&!file[kpaas]:.idea//*&&!file[kpaas]:.github//*` or `!file[kpaas]:.idea//*&&!file[kpaas]:.github//*` for Goland.
7. Click **Apply**.
8. In the **Settings/Preferences** dialog, select **Editor** | **Copyright**.
9. Click **+** Add Copyright icon, and select the scope you just added from the list.
  If you don't see the scope you just added, it may be because your scope does not check **"shared"**
10. From the **Copyright** list, select the profile that just added to link with the scope
11. Click **Apply**.
12. In the **Settings/Preferences** dialog, select **Editor** | **Copyright** | **Formatting**.
13. Configure the formatting options
  * Select **Use line comment**
  * No check **Separator before**
  * No check **Separator after**
  * Select **Before other comments**
  * Check **Add blank line after**

So when you edit the document, the copyright information is automatically added.
For more detailed operations, you can view the official documentation.
[Edit Copyright](https://www.jetbrains.com/help/idea/copyright.html)
[Edit Scope](https://www.jetbrains.com/help/idea/configuring-scopes-and-file-colors.html)

## Community

Since it is still a closed development version, it is not available for the time being. Welcome to leave a message on the issue page

## License

KPaaS source code is licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).
