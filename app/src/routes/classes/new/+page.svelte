<script lang="ts">
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import TextArea from '$lib/components/input/TextArea.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import { pb, processError } from '$lib/pocketbase';
  import { currentUser } from '$lib/stores/user';
  import { defaults, superForm } from 'sveltekit-superforms';
  import { zod } from 'sveltekit-superforms/adapters';
  import { zErrorSchema, zFormSchema } from './models';
  import { goto } from '$app/navigation';

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(
    defaults(zod(zFormSchema)),
    {
      SPA: true,
      validators: zod(zFormSchema),
      async onUpdate({ form }) {
        if (!form.valid) return;

        try {
          await pb.collection('classes').create({ ...form.data, owner: $currentUser!.id });
        } catch (e) {
          return processError(form, e, zErrorSchema);
        }

        goto('/');
      },
    },
  );
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
