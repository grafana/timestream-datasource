import { SelectableValue } from '@grafana/data';

export const sampleQueries: Array<SelectableValue<string>> = [
  {
    label: 'Show databases',
    value: 'SHOW DATABASES',
    description: 'List databases available in your instance',
  },
  {
    label: 'Show tables',
    value: 'SHOW TABLES FROM $__database',
    description: 'List tables in the selected database',
  },
  {
    label: 'Describe table',
    value: 'DESCRIBE $__database.$__table',
    description: 'Describe the selected table',
  },
  {
    label: 'Show measurements',
    value: 'SHOW MEASURES FROM $__database.$__table',
    description: 'List measurements in the selected table',
  },
  {
    label: 'First 10 rows',
    value: 'SELECT * FROM $__database.$__table LIMIT 10',
    description: 'Select the first 10 rows of the selected table',
  },
];
