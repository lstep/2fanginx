// Copyright © 2015 Luc Stepniewski <luc@stepniewski.fr>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

// Thanks to @ariejan, https://ariejan.net/2015/10/12/building-golang-cli-tools-update/

import "fmt"

type version struct {
	Major, Minor, Patch int
	Label               string
	Name                string
}

var Version = version{0, 2, 0, "dev", "Super Duper Thundra"}

var Build string

func (v version) String() string {
	if v.Label != "" {
		return fmt.Sprintf("Version %d.%d.%d-%s \"%s\" Git commit hash: %s", v.Major, v.Minor, v.Patch, v.Label, v.Name, Build)
	} else {
		return fmt.Sprintf("Version %d.%d.%d \"%s\" Git commit hash: %s", v.Major, v.Minor, v.Patch, v.Name, Build)
	}
}
