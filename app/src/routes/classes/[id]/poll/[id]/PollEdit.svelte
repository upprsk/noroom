<script lang="ts">
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import BoolInput from '$lib/components/input/BoolInput.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import { zErrorSchema, zPollSchema } from '$lib/models';
  import { pb, processError } from '$lib/pocketbase';
  import { setMessage, superForm } from 'sveltekit-superforms';
  import { zod } from 'sveltekit-superforms/adapters';
  import { z } from 'zod';

  export let poll: z.infer<typeof zPollSchema>;

  const zFormSchema = zPollSchema.pick({ active: true, title: true });

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(poll, {
    SPA: true,
    validators: zod(zFormSchema),
    async onUpdate({ form }) {
      if (!form.valid) return;

      try {
        await pb.collection('polls').update(poll.id, form.data);
      } catch (e) {
        return processError(form, e, zErrorSchema);
      }

      setMessage(form, 'Salvo');
    },
  });
</script>

<BasicFormCard
  {enhance}
  delayed={$delayed}
  submitting={$submitting}
  _errors={$errors._errors}
  message={$message}
>
  <svelte:fragment slot="title">Editar</svelte:fragment>

  <TextInput
    name="title"
    errors={$errors.title}
    bind:value={$form.title}
    constraints={$constraints.title}
  >
    Titulo
  </TextInput>

  <BoolInput
    name="active"
    errors={$errors.active}
    bind:value={$form.active}
    constraints={$constraints.active}>Ativa</BoolInput
  >

  <svelte:fragment slot="save">Salvar</svelte:fragment>
</BasicFormCard>
