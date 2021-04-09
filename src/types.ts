import { DataQuery, DataSourceJsonData } from '@grafana/data';

export enum Region {
  EntireState = -1,
  SteamboatFlatTops,
  FrontRange,
  VailSummitCounty,
  SawatchRange,
  Aspen,
  Gunnison,
  GrandMesa,
  NorthernSanJuan,
  SouthernSanJuan,
  SangreDeCristo,
}

export interface ZoneQuery extends DataQuery {
  zone?: Region;
}

export const defaultQuery: Partial<ZoneQuery> = {
  zone: Region.EntireState,
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
