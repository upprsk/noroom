<script lang="ts">
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import SelectInput from '$lib/components/input/SelectInput.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import { pb, processError } from '$lib/pocketbase';
  import { defaults, superForm } from 'sveltekit-superforms';
  import { zod } from 'sveltekit-superforms/adapters';
  import { zErrorSchema, zFormSchema } from './models';
  import { goto } from '$app/navigation';

  export let data;

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(
    defaults(zod(zFormSchema)),
    {
      SPA: true,
      validators: zod(zFormSchema),
      async onUpdate({ form }) {
        if (!form.valid) return;

        let res;
        try {
          res = await pb.collection('pods').create(form.data);
        } catch (e) {
          return processError(form, e, zErrorSchema);
        }

        goto(`/pods/${res.id}`);
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
  <svelte:fragment slot="title">Novo Pod</svelte:fragment>
  <svelte:fragment slot="save">Criar</svelte:fragment>

  <TextInput
    name="name"
    errors={$errors.name}
    bind:value={$form.name}
    constraints={$constraints.name}
  >
    Nome
  </TextInput>

  <SelectInput
    name="image"
    errors={$errors.image}
    bind:value={$form.image}
    constraints={$constraints.image}
    options={[{ label: 'Alpine', value: 'alpine' }]}
  >
    Image/Sistema Operacional
  </SelectInput>

  <SelectInput
    name="server"
    errors={$errors.server}
    bind:value={$form.server}
    constraints={$constraints.server}
    options={data.podServers.map((s) => ({ label: s.name, value: s.id }))}
  >
    Servidor
  </SelectInput>
</BasicFormCard>
