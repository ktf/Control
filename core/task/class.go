/*
 * === This file is part of ALICE O² ===
 *
 * Copyright 2018 CERN and copyright holders of ALICE O².
 * Author: Teo Mrnjavac <teo.mrnjavac@cern.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * In applying this license CERN does not waive the privileges and
 * immunities granted to it by virtue of its status as an
 * Intergovernmental Organization or submit itself to any jurisdiction.
 */

package task

import (
	"fmt"
	"github.com/AliceO2Group/Control/common"
	"github.com/AliceO2Group/Control/common/controlmode"
	"github.com/AliceO2Group/Control/core/controlcommands"
	"github.com/AliceO2Group/Control/core/repos"
	"github.com/AliceO2Group/Control/core/task/channel"
	"github.com/AliceO2Group/Control/core/task/constraint"
	"strconv"
)

type TaskClass info

// ↓ We need the roles tree to know *where* to run it and how to *configure* it, but
//   the following information is enough to run the task even with no environment or
//   role info.
type info struct {
	Identifier  taskClassIdentifier		`yaml:"name"`
	Control     struct {
		Mode    controlmode.ControlMode `yaml:"mode"`
	}                                   `yaml:"control"`
	Command     *common.CommandInfo     `yaml:"command"`
	Wants       ResourceWants           `yaml:"wants"`
	Bind        []channel.Inbound       `yaml:"bind"`
	Properties  controlcommands.PropertyMap `yaml:"properties"`
	Constraints []constraint.Constraint `yaml:"constraints"`
}

type taskClassIdentifier struct {
	repo repos.Repo
	Name string
}

func (tcID taskClassIdentifier) String() string {
	if tcID.repo.Revision != "" {
		return fmt.Sprintf("%stasks/%s@%s", tcID.repo.GetIdentifier(), tcID.Name, tcID.repo.Revision)
	} else {
		return fmt.Sprintf("%stasks/%s@master", tcID.repo.GetIdentifier(), tcID.Name)
	}
}

func (tcID *taskClassIdentifier) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	err = unmarshal(&tcID.Name)
	return
}

type ResourceWants struct {
	Cpu     *float64                `yaml:"cpu"`
	Memory  *float64                `yaml:"memory"`
	Ports   Ranges                  `yaml:"ports"`
}

func (rw *ResourceWants) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type _resourceWants struct {
		Cpu     *string                 `yaml:"cpu"`
		Memory  *string                 `yaml:"memory"`
		Ports   *string                 `yaml:"ports"`
	}
	aux := _resourceWants{}
	err = unmarshal(&aux)
	if err != nil {
		return
	}

	if aux.Cpu != nil {
		var cpuCount float64
		cpuCount, err = strconv.ParseFloat(*aux.Cpu, 64)
		if err != nil {
			return
		}
		rw.Cpu = &cpuCount
	}
	if aux.Memory != nil {
		var memCount float64
		memCount, err = strconv.ParseFloat(*aux.Memory, 64)
		if err != nil {
			return
		}
		rw.Memory = &memCount
	}
	if aux.Ports != nil {
		var ranges Ranges
		ranges, err = parsePortRanges(*aux.Ports)
		if err != nil {
			return
		}
		rw.Ports = ranges
	}
	return
}

func (this *info) Equals(other *info) (response bool) {
	if this == nil || other == nil {
		return false
	}
	response = this.Command.Equals(other.Command) &&
		*this.Wants.Cpu == *other.Wants.Cpu &&
		*this.Wants.Memory == *other.Wants.Memory &&
		this.Wants.Ports.Equals(other.Wants.Ports)
	return
}
