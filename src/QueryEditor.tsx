import defaults from 'lodash/defaults';
import React, { PureComponent } from 'react';
import { InlineFormLabel, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, ZoneQuery } from './types';

const ZONES: Array<SelectableValue<number>> = [
  { label: 'Entire State', value: -1 },
  { label: 'Steamboat & Flat Tops', value: 0 },
  { label: 'Front Range', value: 1 },
  { label: 'Vail & Summit County', value: 2 },
  { label: 'Sawatch Range', value: 3 },
  { label: 'Aspen', value: 4 },
  { label: 'Gunnison', value: 5 },
  { label: 'Grand Mesa', value: 6 },
  { label: 'Northern San Juan', value: 7 },
  { label: 'Southern San Juan', value: 8 },
  { label: 'Sangre de Cristo', value: 9 },
  { label: 'Northern Mountains', value: 10 },
  { label: 'Central Mountains', value: 11 },
  { label: 'Southern Mountains', value: 12 },
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
