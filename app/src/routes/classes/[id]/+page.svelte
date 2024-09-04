<script lang="ts">
  import BasicCard from '$lib/components/BasicCard.svelte';
  import { getFileUrl } from '$lib/pocketbase';
  import { currentUser } from '$lib/stores/user';
  import { Marked, type Renderer, type Tokens } from 'marked';
  import { markedHighlight } from 'marked-highlight';
  import Prism from 'prismjs';
  import 'prismjs/components/prism-c';
  import 'prismjs/components/prism-cpp';
  import 'prismjs/components/prism-json';
  import 'prismjs/themes/prism-tomorrow.css';
  import DOMPurify from 'dompurify';
  import { onMount } from 'svelte';
  import BasicAvatarNoLink from '$lib/components/BasicAvatarNoLink.svelte';

  Prism.manual = true;

  const renderer = {
    html(this: Renderer, token: Tokens.HTML | Tokens.Tag): string {
      const re = /^<!--\s+end_slide\s+-->/;
      if (!re.test(token.text)) return token.text;

      return `<div class="divider"></div>`;
    },
    image(this: Renderer, token: Tokens.Image): string {
      const m = token.href.match(/^\.\/(.+?)\.\w+$/);
      if (!m) return `<img alt="${token.text}" href="${token.href}"></img>`;

      const [, filename] = m;

      const s = data.klass.attachments
        .map((a) => a.match(/^(.+?)_\w{10}\.\w+$/))
        .filter((a) => a !== null)
        .find(([, a]) => a === filename);
      if (!s) return `<img alt="${token.text}" href="${token.href}"></img>`;

      return `<img alt="${token.text}" src="${getFileUrl(data.klass, s[0])}"></img>`;
    },
  };

  const marked = new Marked(
    { renderer },
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
      },
    }),
  );

  export let data;

  let sanitize: undefined | ((s: string) => void);
  onMount(async () => {
    sanitize = DOMPurify.sanitize;
  });
</script>

<BasicCard>
  <svelte:fragment slot="title">{data.klass.title}</svelte:fragment>

  {#if data.classPolls.length > 0}
    <div class="flex gap-2">
      {#each data.classPolls as poll (poll.id)}
        <a href="poll/{poll.id}" class="btn" class:btn-primary={poll.active}>
          {poll.title}
          <div class="badge">{poll.expand.pollAnswers_via_poll?.length ?? 0}</div>
        </a>
      {/each}

      {#if data.user.role === 'editor'}
        <a href="poll/new" class="btn">
          Novo <div class="badge">+</div>
        </a>
      {/if}
    </div>
  {/if}

  <div class="prose">
    <!-- {@html htmlContent} -->
    {#await marked.parse(data.klass.content)}
      <!-- promise is pending -->
      pending...
    {:then value}
      {#if sanitize}
        {@html DOMPurify.sanitize(value)}
      {/if}
    {/await}
  </div>

  {#if data.klass.attachments.length > 0}
    <div class="divider"></div>

    <div class="flex w-full flex-wrap gap-2">
      {#each data.klass.attachments as attach}
        <a href={getFileUrl(data.klass, attach)} class="btn btn-sm">
          <svg class="h-4 w-4 fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256"
            ><path
              d="M213.66,82.34l-56-56A8,8,0,0,0,152,24H56A16,16,0,0,0,40,40V216a16,16,0,0,0,16,16H200a16,16,0,0,0,16-16V88A8,8,0,0,0,213.66,82.34ZM160,51.31,188.69,80H160ZM200,216H56V40h88V88a8,8,0,0,0,8,8h48V216Z"
            ></path></svg
          >
          {attach}
        </a>
      {/each}
    </div>
  {/if}

  {#if data.klass.expand?.classPresenceEntries_via_class && data.klass.expand?.classPresenceEntries_via_class.length > 0}
    <div class="divider"></div>

    <div class="flex w-full flex-wrap items-center">
      {#if data.klass?.expand?.classPresenceEntries_via_class}
        <div class="avatar-group -space-x-6 p-2 rtl:space-x-reverse">
          {#each data.klass.expand.classPresenceEntries_via_class as presence}
            <BasicAvatarNoLink user={presence.expand.user} />
          {/each}

          <div class="avatar placeholder">
            <div class="w-10 bg-neutral text-neutral-content">
              <span>{data.klass.expand.classPresenceEntries_via_class.length}</span>
            </div>
          </div>
        </div>

        {#if data.klass.expand.classPresenceEntries_via_class.find((p) => p.user === $currentUser?.id)}
          <div class="badge badge-success">Presente</div>
        {/if}
      {/if}
    </div>
  {/if}

  <div class="card-actions justify-end">
    <button type="button" class="btn" on:click={() => history.back()}>voltar</button>
    {#if data.klass.live}
      <a href="presence" class="btn">presenca</a>
    {/if}
    {#if $currentUser?.id === data.klass.owner || $currentUser?.role === 'editor'}
      <a href="edit" class="btn btn-primary">editar</a>
    {/if}
  </div>
</BasicCard>
