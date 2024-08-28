<script lang="ts">
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import TextArea from '$lib/components/input/TextArea.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import { superForm } from 'sveltekit-superforms';

  export let data;

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(data.form);
</script>

<BasicFormCard
  {enhance}
  delayed={$delayed}
  submitting={$submitting}
  _errors={$errors._errors}
  message={$message}
>
  <svelte:fragment slot="title">Nova aula</svelte:fragment>
  <svelte:fragment slot="save">Criar</svelte:fragment>

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
</BasicFormCard>
