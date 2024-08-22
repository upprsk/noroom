<script lang="ts">
  import { browser } from '$app/environment';
  import BasicCard from '$lib/components/BasicCard.svelte';
  import DOMPurify from 'dompurify';
  import { Marked } from 'marked';
  import { markedHighlight } from 'marked-highlight';
  import Prism from 'prismjs';
  import 'prismjs/components/prism-c';
  import 'prismjs/components/prism-cpp';
  import 'prismjs/components/prism-json';
  import 'prismjs/themes/prism-tomorrow.css';

  Prism.manual = true;

  const sanitize = browser
    ? DOMPurify.sanitize
    : (async () => {
        const { JSDOM } = await import('jsdom');
        const window = new JSDOM('').window;
        return DOMPurify(window).sanitize;
      })();

  const marked = new Marked(
    markedHighlight({
      langPrefix: 'language-',
      highlight(code, lang) {
        let grammar = Prism.languages.clike;
        switch (lang) {
          case 'c':
            grammar = Prism.languages.c;
            break;
          case 'cpp':
            grammar = Prism.languages.cpp;
            break;
          case 'json':
            grammar = Prism.languages.json;
            break;
        }

        return Prism.highlight(code, grammar, lang);
        // const language = hljs.getLanguage(lang) ? lang : 'plaintext';
        // return hljs.highlight(code, { language }).value;
      },
    }),
  );

  export let data;
</script>

<BasicCard>
  <svelte:fragment slot="title">{data.klass.title}</svelte:fragment>

  <div class="prose">
    <!-- {@html htmlContent} -->
    {#await marked.parse(data.klass.content)}
      <!-- promise is pending -->
    {:then value}
      {#await sanitize then s}
        {@html s(value)}
      {/await}
    {/await}
  </div>
</BasicCard>
