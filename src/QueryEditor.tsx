import React, { PureComponent } from 'react';
import { Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { MyDataSourceOptions, ZoneQuery } from './types';

const ZONES: Array<SelectableValue<string>> = [
  { label: 'Entire State', value: 'entire_statea' },
  { label: 'Steamboat & Flat Tops', value: 'steamboat_flat_tops' },
  { label: 'Front Range', value: 'front_range' },
  { label: 'Vail & Summit County', value: 'vail_summit_county' },
  { label: 'Sawatch Range', value: 'sawatch_range' },
  { label: 'Aspen', value: 'aspen' },
  { label: 'Gunnison', value: 'gunnison' },
  { label: 'Grand Mesa', value: 'grand_mesa' },
  { label: 'Northern San Juan', value: 'northern_san_juan' },
  { label: 'Southern San Juan', value: 'southern_san_juan' },
  { label: 'Sangre de Cristo', value: 'sangre_de_criso' },
  { label: 'Northern Mountains', value: 'northern_mountains' },
  { label: 'Central Mountains', value: 'central_mountians' },
  { label: 'Southern Mountains', value: 'southern_mountains' },
];

const DEFAULT_ZONE = {
  label: 'Entire State',
  value: 'entire_state',
};

type Props = QueryEditorProps<DataSource, ZoneQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onRegionChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, zone: value.value });
    onRunQuery();
  };

  render() {
    return (
      <div className="gf-form">
        <Select width={20} options={ZONES} value={DEFAULT_ZONE} onChange={this.onRegionChange} />
      </div>
    );
  }
}
