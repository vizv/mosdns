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
	"github.com/IrineSistiana/mosdns/v5/pkg/matcher/domain"
	"github.com/v2fly/v2ray-core/v4/app/router"
)

func loadFromGeosite[T any](m *domain.MixMatcher[struct{}], geositeList *router.GeoSiteList, cs map[string]bool) error {
	for _, entry := range geositeList.Entry {
		if !cs[entry.CountryCode] {
			continue
		}

		for _, dom := range entry.Domain {
			var pattern = dom.Value

			switch dom.Type {
			case router.Domain_Full:
				pattern = domain.MatcherFull + ":" + pattern
			case router.Domain_Domain:
				pattern = domain.MatcherDomain + ":" + pattern
			case router.Domain_Regex:
				pattern = domain.MatcherRegexp + ":" + pattern
			case router.Domain_Plain:
				pattern = domain.MatcherKeyword + ":" + pattern
			default:
				continue
			}

			if err := m.Add(pattern, struct{}{}); err != nil {
				return err
			}
		}
	}

	return nil
}
