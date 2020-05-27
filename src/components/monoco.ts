import { monaco } from '@monaco-editor/react';

monaco
  .init()
  .then(monaco => {
    console.log('init!!!', monaco);

    if (false) {
      // Register a completion item provider for the new language
      monaco.languages.registerCompletionItemProvider('sql', {
        provideCompletionItems: (model: any, position: any, context: any) => {
          console.log('Complete???', model, position, context);
          var suggestions: any[] = [
            {
              label: 'simpleText',
              kind: monaco.languages.CompletionItemKind.Text,
              insertText: 'simpleText',
            },
            {
              label: 'testing',
              kind: monaco.languages.CompletionItemKind.Keyword,
              insertText: 'testing(${1:condition})',
              insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
              label: 'ifelse',
              kind: monaco.languages.CompletionItemKind.Snippet,
              insertText: ['if (${1:condition}) {', '\t$0', '} else {', '\t', '}'].join('\n'),
              insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
              documentation: 'If-Else Statement',
            },
          ];

          return { suggestions }; //   ProviderResul<CompletionList>
        },
      });
    }

    console.log('done', monaco);
  })
  .catch(error => console.error('An error occurred during initialization of Monaco: ', error));
