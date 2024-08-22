<script lang="ts">
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import TextArea from '$lib/components/input/TextArea.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import { superForm } from 'sveltekit-superforms';
  import { applyAction, enhance as enhanceForm } from '$app/forms';

  export let data;

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(
    data.form,
    { invalidateAll: 'force' },
  );
</script>

<BasicFormCard
  {enhance}
  delayed={$delayed}
  submitting={$submitting}
  _errors={$errors._errors}
  message={$message}
  action="?/save"
>
  <svelte:fragment slot="title">Editar {$form.title}</svelte:fragment>
  <svelte:fragment slot="save">salvar</svelte:fragment>

  <TextInput
    name="title"
    errors={$errors.title}
    bind:value={$form.title}
    constraints={$constraints.title}
  >
    Title
  </TextInput>

  <TextArea
    name="content"
    errors={$errors.content}
    bind:value={$form.content}
    constraints={$constraints.content}
  >
    Content
  </TextArea>

  <div class="flex w-full gap-2">
    {#each $form.attachments as attach}
      <button
        class="btn btn-sm"
        type="button"
        on:click={() => ($form.attachments = $form.attachments.filter((it) => it !== attach))}
      >
        <svg class="w-4 h-4 fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256"
          ><path
            d="M213.66,82.34l-56-56A8,8,0,0,0,152,24H56A16,16,0,0,0,40,40V216a16,16,0,0,0,16,16H200a16,16,0,0,0,16-16V88A8,8,0,0,0,213.66,82.34ZM160,51.31,188.69,80H160ZM200,216H56V40h88V88a8,8,0,0,0,8,8h48V216Z"
          ></path></svg
        >
        {attach}
      </button>
    {/each}
  </div>

  <svelte:fragment slot="actions">
    <button type="button" class="btn" on:click={() => history.back()}>voltar</button>
    <a href="attach-files" class="btn">anexar</a>
    <form
      method="POST"
      action="?/remove"
      use:enhanceForm={() =>
        async ({ result }) => {
          await applyAction(result);
        }}
    >
      <button class="btn btn-warning">remover</button>
    </form>
  </svelte:fragment>
</BasicFormCard>
