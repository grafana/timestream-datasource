import React, { PureComponent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import Editor from '@monaco-editor/react';
import { DataSource } from '../DataSource';
import { TimestreamQuery, TimestreamOptions } from '../types';
import { config } from '@grafana/runtime';

import './monoco';

type Props = QueryEditorProps<DataSource, TimestreamQuery, TimestreamOptions>;

/**
 *
 * https://github.com/influxdata/influxdb/blob/master/ui/src/shared/components/FluxMonacoEditor.tsx
 * https://microsoft.github.io/monaco-editor/playground.html#extending-language-services-completion-provider-example
 *
 */

export class QueryEditor extends PureComponent<Props> {
  getEditorValue: any | undefined;

  onRawQueryChange = () => {
    this.props.onChange({
      ...this.props.query,
      rawQuery: this.getEditorValue(),
    });
    this.props.onRunQuery();
  };

  onEditorDidMount = (getEditorValue: any) => {
    this.getEditorValue = getEditorValue;
  };

  render() {
    const { query } = this.props;

    return (
      <>
        <div onBlur={this.onRawQueryChange}>
          <Editor
            height={'250px'}
            language="sql"
            value={query.rawQuery}
            editorDidMount={this.onEditorDidMount}
            theme={config.theme.isDark ? 'dark' : 'light'}
            options={{
              wordWrap: 'off',
              codeLens: false, // too small to bother
              minimap: {
                enabled: false,
                renderCharacters: false,
              },
              lineNumbersMinChars: 4,
              lineDecorationsWidth: 0,
              overviewRulerBorder: false,
              automaticLayout: true,
            }}
          />
        </div>
      </>
    );
  }
}
