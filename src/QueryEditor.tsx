import defaults from 'lodash/defaults';
import React, { PureComponent } from 'react';
import { InlineFormLabel, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, ZoneQuery } from './types';

const ZONES: Array<SelectableValue<string>> = [
  { label: 'Entire State', value: 'entire-state' },
  { label: 'Steamboat & Flat Tops', value: 'steamboat-flat-tops' },
  { label: 'Front Range', value: 'front-range' },
  { label: 'Vail & Summit County', value: 'vail-summit-county' },
  { label: 'Sawatch Range', value: 'sawatch' },
  { label: 'Aspen', value: 'aspen' },
  { label: 'Gunnison', value: 'gunnison' },
  { label: 'Grand Mesa', value: 'grand-mesa' },
  { label: 'Northern San Juan', value: 'north-san-juan' },
  { label: 'Southern San Juan', value: 'south-san-juan' },
  { label: 'Sangre de Cristo', value: 'sangre-de-cristo' },
];

type Props = QueryEditorProps<DataSource, ZoneQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onRegionChange = (value: SelectableValue<string>) => {
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
