import { LanguageDefinition } from '@grafana/experimental';
import { conf, language } from './language';

const timestreamLanguageDefinition: LanguageDefinition & { id: string } = {
  id: 'timestream',
  // TODO: Load language using code splitting instead: loader: () => import('./language'),
  loader: () => Promise.resolve({ conf, language }),
};

export default timestreamLanguageDefinition;
