/*
 * Copyright (C) 2023, VizV
 *
 * This file is part of mosdns.
 *
 * mosdns is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mosdns is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package v2ray_geosite

import (
	"fmt"
	"os"
	"strings"

	"github.com/IrineSistiana/mosdns/v5/coremain"
	"github.com/IrineSistiana/mosdns/v5/pkg/matcher/domain"
	"github.com/IrineSistiana/mosdns/v5/plugin/data_provider"
	"github.com/v2fly/v2ray-core/v4/app/router"
	"google.golang.org/protobuf/proto"
)

const PluginType = "v2ray_geosite"

var geositeListFiles = map[string]*router.GeoSiteList{}

func init() {
	coremain.RegNewPluginFunc(PluginType, Init, func() any { return new(Args) })
}

func Init(bp *coremain.BP, args any) (any, error) {
	m, err := NewV2rayGeosite(bp, args.(*Args))
	if err != nil {
		return nil, err
	}
	return m, nil
}

type Args struct {
	Sets  []string `yaml:"sets"`
	Files []string `yaml:"files"`
	Codes []string `yaml:"codes"`
}

var _ data_provider.DomainMatcherProvider = (*V2rayGeosite)(nil)

type V2rayGeosite struct {
	mg []domain.Matcher[struct{}]
}

func (d *V2rayGeosite) GetDomainMatcher() domain.Matcher[struct{}] {
	return MatcherGroup(d.mg)
}

// NewV2rayGeosite inits a V2rayGeosite from given args.
func NewV2rayGeosite(bp *coremain.BP, args *Args) (*V2rayGeosite, error) {
	v2gs := &V2rayGeosite{}

	m := domain.NewDomainMixMatcher()

	cs := map[string]bool{}
	for _, code := range args.Codes {
		cs[strings.ToUpper(code)] = true
	}

	if err := LoadFiles(args.Files, cs, m); err != nil {
		return nil, err
	}
	if m.Len() > 0 {
		v2gs.mg = append(v2gs.mg, m)
	}

	for _, tag := range args.Sets {
		provider, _ := bp.M().GetPlugin(tag).(data_provider.DomainMatcherProvider)
		if provider == nil {
			return nil, fmt.Errorf("%s is not a DomainMatcherProvider", tag)
		}
		m := provider.GetDomainMatcher()
		v2gs.mg = append(v2gs.mg, m)
	}
	return v2gs, nil
}

func LoadFiles(fs []string, cs map[string]bool, m *domain.MixMatcher[struct{}]) error {
	for i, f := range fs {
		if err := LoadFile(f, cs, m); err != nil {
			return fmt.Errorf("failed to load file #%d %s, %w", i, f, err)
		}
	}
	return nil
}

func LoadFile(f string, cs map[string]bool, m *domain.MixMatcher[struct{}]) error {
	if len(f) > 0 {
		var geositeList = geositeListFiles[f]
		if geositeList == nil {
			geositeList = &router.GeoSiteList{}
			data, err := os.ReadFile(f)
			if err != nil {
				return err
			}
			if err := proto.Unmarshal(data, geositeList); err != nil {
				return err
			}
			geositeListFiles[f] = geositeList
		}

		if err := loadFromGeosite[struct{}](m, geositeList, cs); err != nil {
			return err
		}
	}
	return nil
}
