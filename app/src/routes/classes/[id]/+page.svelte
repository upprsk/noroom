<script lang="ts">
  import { browser } from '$app/environment';
  import BasicCard from '$lib/components/BasicCard.svelte';
  import { getFileUrl } from '$lib/pocketbase.js';
  import { currentUser } from '$lib/stores/user';
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

  {#if data.klass.attachments.length > 0}
    <div class="divider"></div>

    <div class="flex w-full gap-2">
      {#each data.klass.attachments as attach}
        <a href={getFileUrl(data.klass, attach)} class="btn btn-sm">
          <svg class="w-4 h-4 fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256"
            ><path
              d="M213.66,82.34l-56-56A8,8,0,0,0,152,24H56A16,16,0,0,0,40,40V216a16,16,0,0,0,16,16H200a16,16,0,0,0,16-16V88A8,8,0,0,0,213.66,82.34ZM160,51.31,188.69,80H160ZM200,216H56V40h88V88a8,8,0,0,0,8,8h48V216Z"
            ></path></svg
          >
          {attach}
        </a>
      {/each}
    </div>
  {/if}

  <div class="card-actions justify-end">
    {#if $currentUser?.id === data.klass.owner || $currentUser?.role === 'editor'}
      <a href="{data.klass.id}/edit" class="btn btn-primary">editar</a>
    {/if}
  </div>
</BasicCard>
