import { SelectableValue } from '@grafana/data';

export const standardRegions: Array<SelectableValue<string>> = ['us-east-1', 'us-east-2', 'us-west-2', 'eu-west-2'].map(
  r => {
    return { value: r, label: r };
  }
);
