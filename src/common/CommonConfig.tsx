import React, { PureComponent } from 'react';
import { InlineFormLabel, LegacyForms, Button } from '@grafana/ui';
const { Select, Input } = LegacyForms;
import {
  DataSourcePluginOptionsEditorProps,
  onUpdateDatasourceJsonDataOptionSelect,
  onUpdateDatasourceResetOption,
  onUpdateDatasourceJsonDataOption,
  onUpdateDatasourceSecureJsonDataOption,
  SelectableValue,
} from '@grafana/data';

import {
  AwsDataSourceJsonData,
  AwsDataSourceSecureJsonData,
  awsAuthProviderOptions,
  AwsAuthType,
  standardRegions,
} from './types';

export interface Props extends DataSourcePluginOptionsEditorProps<AwsDataSourceJsonData, AwsDataSourceSecureJsonData> {
  loadRegions?: () => Promise<string[]>;
}

export interface State {
  regions: Array<SelectableValue<string>>;
}

export default class CommonConfig extends PureComponent<Props, State> {
  state: State = {
    regions: standardRegions.map(r => {
      return { value: r, label: r };
    }),
  };

  // loadRegionsPromise: CancelablePromis<any> | null = null;

  // componentDidMount() {
  //   this.loadRegionsPromise = makePromiseCancelable(this.loadRegions());
  //   this.loadRegionsPromise.promise.catch(({ isCanceled }) => {
  //     if (isCanceled) {
  //       console.warn('Cloud Watch ConfigEditor has unmounted, intialization was canceled');
  //     }
  //   });
  // }

  // componentWillUnmount() {
  //   if (this.loadRegionsPromise) {
  //     this.loadRegionsPromise.cancel();
  //   }
  // }

  render() {
    const { regions } = this.state;
    const { options } = this.props;
    const jsonData = (options.jsonData || {}) as AwsDataSourceJsonData;
    const secureJsonData = (options.secureJsonData || {}) as AwsDataSourceSecureJsonData;

    const accessKeyConfigured = options?.secureJsonFields?.accessKey === true;
    const secretKeyConfigured = options?.secureJsonFields?.secretKey === true;

    const widthKey = 'width-12';
    const widthVal = 'width-30';

    return (
      <>
        <div className="gf-form-inline">
          <div className="gf-form">
            <InlineFormLabel className={widthKey}>Auth Provider</InlineFormLabel>
            <Select
              className={widthVal}
              value={awsAuthProviderOptions.find(authProvider => authProvider.value === jsonData.authType)}
              options={awsAuthProviderOptions}
              defaultValue={jsonData.authType}
              onChange={option => {
                // Remove it when switching away from ARN
                if (jsonData.authType === AwsAuthType.ARN && option.value !== AwsAuthType.ARN) {
                  delete this.props.options.jsonData.assumeRoleArn;
                }
                onUpdateDatasourceJsonDataOptionSelect(this.props, 'authType')(option);
              }}
            />
          </div>
        </div>
        {jsonData.authType === AwsAuthType.Credentials && (
          <div className="gf-form-inline">
            <div className="gf-form">
              <InlineFormLabel
                className={widthKey}
                tooltip="Credentials profile name, as specified in ~/.aws/credentials, leave blank for default."
              >
                Credentials Profile Name
              </InlineFormLabel>
              <div className={widthVal}>
                <Input
                  className={widthVal}
                  placeholder="~/.aws/credentials"
                  value={jsonData.profile}
                  onChange={onUpdateDatasourceJsonDataOption(this.props, 'profile')}
                />
              </div>
            </div>
          </div>
        )}
        {jsonData.authType === AwsAuthType.Keys && (
          <>
            <div className="gf-form-inline">
              <div className="gf-form">
                <InlineFormLabel className={widthKey}>Access Key ID</InlineFormLabel>
                {accessKeyConfigured ? (
                  <Input className={widthVal} placeholder="saved" disabled={true} />
                ) : (
                  <Input
                    className={widthVal}
                    value={secureJsonData.accessKey ?? ''}
                    placeholder="Access Key ID"
                    required
                    onChange={onUpdateDatasourceSecureJsonDataOption(this.props, 'accessKey')}
                  />
                )}
              </div>
              {accessKeyConfigured && (
                <div className="gf-form">
                  <div className="max-width-30 gf-form-inline">
                    <Button
                      variant="secondary"
                      type="button"
                      onClick={onUpdateDatasourceResetOption(this.props as any, 'accessKey')}
                    >
                      reset
                    </Button>
                  </div>
                </div>
              )}
            </div>

            <div className="gf-form-inline">
              <div className="gf-form">
                <InlineFormLabel className={widthKey}>Secret Key</InlineFormLabel>
                {secretKeyConfigured ? (
                  <Input className={widthVal} placeholder="saved" disabled={true} />
                ) : (
                  <Input
                    className={widthVal}
                    value={secureJsonData.secretKey ?? ''}
                    placeholder="Secret Key"
                    required
                    onChange={onUpdateDatasourceSecureJsonDataOption(this.props, 'secretKey')}
                  />
                )}
              </div>
              {secretKeyConfigured && (
                <div className="gf-form">
                  <div className="max-width-30 gf-form-inline">
                    <Button
                      variant="secondary"
                      type="button"
                      onClick={onUpdateDatasourceResetOption(this.props as any, 'secretKey')}
                    >
                      reset
                    </Button>
                  </div>
                </div>
              )}
            </div>
          </>
        )}
        {jsonData.authType === AwsAuthType.ARN && (
          <div className="gf-form-inline">
            <div className="gf-form">
              <InlineFormLabel className={widthKey} tooltip="ARN of Assume Role">
                Assume Role ARN
              </InlineFormLabel>
              <div className={widthVal}>
                <Input
                  className={widthVal}
                  placeholder="arn:aws:iam:*"
                  value={options.jsonData.assumeRoleArn || ''}
                  onChange={onUpdateDatasourceJsonDataOption(this.props, 'assumeRoleArn')}
                />
              </div>
            </div>
          </div>
        )}
        <div className="gf-form-inline">
          <div className="gf-form">
            <InlineFormLabel
              className={widthKey}
              tooltip="Specify the region, such as for US West (Oregon) use ` us-west-2 ` as the region."
            >
              Default Region
            </InlineFormLabel>
            <Select
              className={widthVal}
              value={regions.find(region => region.value === jsonData.defaultRegion)}
              options={regions}
              defaultValue={jsonData.defaultRegion}
              onChange={onUpdateDatasourceJsonDataOptionSelect(this.props, 'defaultRegion')}
            />
          </div>
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <InlineFormLabel className={widthKey} tooltip="Override the default service endpoint">
              Endpoint (optional)
            </InlineFormLabel>
            <div className={widthVal}>
              <Input
                className={widthVal}
                placeholder="https://query-{cell}.timestream.{region}.amazonaws.com"
                value={jsonData.endpoint}
                onChange={onUpdateDatasourceJsonDataOption(this.props, 'endpoint')}
              />
            </div>
          </div>
        </div>
      </>
    );
  }
}
