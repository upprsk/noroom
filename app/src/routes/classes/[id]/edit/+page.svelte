<script lang="ts">
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import TextArea from '$lib/components/input/TextArea.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import { superForm } from 'sveltekit-superforms';
  import { enhance as enhanceForm } from '$app/forms';
  import ListAndRemoveAttachments from './ListAndRemoveAttachments.svelte';
  import { currentUser } from '$lib/stores/user';

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

  <ListAndRemoveAttachments klass={data.klass} />

  <svelte:fragment slot="actions">
    <button type="button" class="btn" on:click={() => history.back()}>voltar</button>
    <a href="attach-files" class="btn">anexar</a>
    {#if $currentUser?.role === 'editor'}
      <form method="POST" action="?/remove" use:enhanceForm>
        <button class="btn btn-warning">remover</button>
      </form>
    {/if}
  </svelte:fragment>
</BasicFormCard>
