import { MacroType } from '@grafana/experimental';

export const TABLE_MACRO = '$__table';
export const DATABASE_MACRO = '$__database';

export const MACROS = [
  {
    id: '$__timeFilter',
    name: '$__timeFilter',
    text: '$__timeFilter',
    args: [],
    type: MacroType.Filter,
    description: 'Will be replaced by an expression that limits the time to the dashboard range.',
  },
  {
    id: '$__timeFrom',
    name: '$__timeFrom',
    text: '$__timeFrom',
    args: [],
    type: MacroType.Filter,
    description: 'Will be replaced by the number in milliseconds at the start of the dashboard range',
  },
  {
    id: '$__timeTo',
    name: '$__timeTo',
    text: '$__timeTo',
    args: [],
    type: MacroType.Filter,
    description: 'Will be replaced by the number in milliseconds at the end of the dashboard range.',
  },
  {
    id: '$__interval_ms',
    name: '$__interval_ms',
    text: '$__interval_ms',
    args: [],
    type: MacroType.Filter,
    description:
      'Will be replaced by a number in time format that represents the amount of time a single pixel in the graph should cover.',
  },
  {
    id: '$__interval_raw_ms',
    name: '$__interval_raw_ms',
    text: '$__interval_raw_ms',
    args: [],
    type: MacroType.Filter,
    description:
      'Will be replaced by the number in milliseconds that represents the amount of time a single pixel in the graph should cover.',
  },
  {
    id: DATABASE_MACRO,
    name: DATABASE_MACRO,
    text: DATABASE_MACRO,
    args: [],
    type: MacroType.Table,
    description: 'Will be replaced by the query database.',
  },
  {
    id: TABLE_MACRO,
    name: TABLE_MACRO,
    text: TABLE_MACRO,
    args: [],
    type: MacroType.Table,
    description: 'Will be replaced by the query table.',
  },
  {
    id: '$__measure',
    name: '$__measure',
    text: '$__measure',
    args: [],
    type: MacroType.Column,
    description: 'Will be replaced by the query measure.',
  },
];
