<script lang="ts">
  import { goto } from '$app/navigation';
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import SelectInput from '$lib/components/input/SelectInput.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import { zErrorSchema, zPollSchema } from '$lib/models.js';
  import { pb, processError } from '$lib/pocketbase.js';
  import { superForm, defaults } from 'sveltekit-superforms';
  import { zod } from 'sveltekit-superforms/adapters';

  export let data;

  const zFormSchema = zPollSchema;

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(
    defaults(zod(zFormSchema)),
    {
      SPA: true,
      validators: zod(zFormSchema),
      async onUpdate({ form }) {
        if (!form.valid) return;

        let res;
        try {
          res = await pb.collection('polls').create({ ...form.data, class: data.klass.id });
        } catch (e) {
          return processError(form, e, zErrorSchema);
        }

        goto(`/classes/${data.klass.id}/poll/${res.id}`);
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
  <svelte:fragment slot="title">Nova Poll</svelte:fragment>
  <svelte:fragment slot="save">Criar</svelte:fragment>

  <TextInput
    name="title"
    errors={$errors.title}
    bind:value={$form.title}
    constraints={$constraints.title}
  >
    Titulo
  </TextInput>

  <SelectInput
    name="expects"
    errors={$errors.expects}
    bind:value={$form.expects}
    constraints={$constraints.expects}
    options={zFormSchema.shape.expects.options}
  >
    Tipo
  </SelectInput>

  {#if $form.expects === 'string'}
    <p>Utilizando "string"</p>
  {:else if $form.expects === 'number'}
    <p>Utilizando "number"</p>
  {:else if $form.expects === 'option'}
    <TextInput
      name="options"
      errors={$errors.options}
      bind:value={$form.options}
      constraints={$constraints.options}
    >
      Opcoes
      <span class="label-text" slot="extra-top">separar por ","</span>
    </TextInput>
  {:else if $form.expects === 'multi'}
    <TextInput
      name="options"
      errors={$errors.options}
      bind:value={$form.options}
      constraints={$constraints.options}
    >
      Opcoes
      <span class="label-text" slot="extra-top">separar por ","</span>
    </TextInput>
  {/if}
</BasicFormCard>
