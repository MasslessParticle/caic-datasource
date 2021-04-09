import defaults from 'lodash/defaults';
import React, { PureComponent } from 'react';
import { InlineFormLabel, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, Region, ZoneQuery } from './types';

const ZONES: Array<SelectableValue<number>> = [
  { label: 'Entire State', value: Region.EntireState },
  { label: 'Steamboat & Flat Tops', value: Region.SteamboatFlatTops },
  { label: 'Front Range', value: Region.FrontRange },
  { label: 'Vail & Summit County', value: Region.VailSummitCounty },
  { label: 'Sawatch Range', value: Region.SawatchRange },
  { label: 'Aspen', value: Region.Aspen },
  { label: 'Gunnison', value: Region.Gunnison },
  { label: 'Grand Mesa', value: Region.GrandMesa },
  { label: 'Northern San Juan', value: Region.NorthernSanJuan },
  { label: 'Southern San Juan', value: Region.SouthernSanJuan },
  { label: 'Sangre de Cristo', value: Region.SangreDeCristo },
];

type Props = QueryEditorProps<DataSource, ZoneQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onRegionChange = (value: SelectableValue<number>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, zone: value.value });
    onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { zone } = query;

    return (
      <div className="gf-form">
        <div className="gf-form-inline">
          <InlineFormLabel width={12} className="zone-label" tooltip="select a geographic zone">
            Select a Geographic Zone
          </InlineFormLabel>
          <Select width={30} options={ZONES} value={zone} onChange={this.onRegionChange} />
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
      </div>
    );
  }
}
