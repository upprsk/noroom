<script lang="ts">
  import { page } from '$app/stores';
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import TextArea from '$lib/components/input/TextArea.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import { pb, processError } from '$lib/pocketbase';
  import { currentUser } from '$lib/stores/user';
  import { setMessage, superForm } from 'sveltekit-superforms';
  import { zod } from 'sveltekit-superforms/adapters';
  import ListAndRemoveAttachments from './ListAndRemoveAttachments.svelte';
  import { zErrorSchema, zSaveSchema } from './models';
  import { goto, invalidate } from '$app/navigation';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import RemoveWithConfirm from './RemoveWithConfirm.svelte';

  export let data;

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(
    data.form,
    {
      SPA: true,
      validators: zod(zSaveSchema),
      invalidateAll: 'force',
      async onUpdate({ form }) {
        if (!form.valid) return;

        try {
          await pb
            .collection('classes')
            .update($page.params.id, { ...form.data, owner: undefined });
        } catch (e) {
          return processError(form, e, zErrorSchema);
        }

        return setMessage(form, 'Salvo');
      },
    },
  );

  let removeErrors: string[] | undefined;
  const remove = async () => {
    try {
      removeErrors = undefined;

      await pb.collection('classes').delete($page.params.id);
    } catch (e) {
      console.error(e);
      removeErrors = [`${e}`];
    }

    goto('/');
  };

  const removeAttachment = async (attachment: string) => {
    try {
      removeErrors = undefined;

      await pb.collection('classes').update($page.params.id, { 'attachments-': attachment });
      invalidate('app:klass');
    } catch (e) {
      console.error(e);
      removeErrors = [`${e}`];
    }
  };
</script>

<BasicFormCard
  {enhance}
  delayed={$delayed}
  submitting={$submitting}
  _errors={$errors._errors}
  message={$message}
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

  <ListAndRemoveAttachments klass={data.klass} on:remove={(e) => removeAttachment(e.detail)} />

  <ErrorAlert errors={removeErrors} />

  <svelte:fragment slot="actions">
    <button type="button" class="btn" on:click={() => history.back()}>voltar</button>
    <a href="attach-files" class="btn">anexar</a>
    {#if $currentUser?.role === 'editor'}
      <!-- <button type="button" class="btn btn-warning" on:click={remove}>remover</button> -->
      <RemoveWithConfirm on:remove={remove} />
    {/if}
  </svelte:fragment>
</BasicFormCard>
