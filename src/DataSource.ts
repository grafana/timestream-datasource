import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { TimestreamQuery, TimestreamOptions } from './types';

export class DataSource extends DataSourceWithBackend<TimestreamQuery, TimestreamOptions> {
  templateSrv: any;

  constructor(instanceSettings: DataSourceInstanceSettings<TimestreamOptions>) {
    super(instanceSettings);

    this.templateSrv = getTemplateSrv();
  }

  applyTemplateVariables(query: TimestreamQuery): TimestreamQuery {
    return {
      ...query,
      rawQuery: this.templateSrv.replace(query.rawQuery),
    };
  }
}
