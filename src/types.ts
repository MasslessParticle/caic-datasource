import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface ZoneQuery extends DataQuery {
  zone?: number;
}

export const defaultQuery: Partial<ZoneQuery> = {
  zone: -1,
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apiKey?: string;
}
