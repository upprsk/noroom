<script lang="ts">
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import NumberInput from '$lib/components/input/NumberInput.svelte';
  import SelectInput from '$lib/components/input/SelectInput.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import {
    zMakeErrorDataSchema,
    zPollAnswerSchema,
    type zPollSchema,
    type zUserSchema,
  } from '$lib/models';
  import { pb, processError } from '$lib/pocketbase';
  import { defaults, superForm } from 'sveltekit-superforms';
  import { zod } from 'sveltekit-superforms/adapters';
  import { z } from 'zod';

  export let user: z.infer<typeof zUserSchema>;
  export let poll: z.infer<typeof zPollSchema>;

  const zFormSchema =
    poll.expects === 'string'
      ? zPollAnswerSchema
          .omit({ answer: true })
          .extend({ answer: z.string(), expects: z.literal('string').default('string') })
      : poll.expects === 'option'
        ? zPollAnswerSchema
            .omit({ answer: true })
            .extend({ answer: z.string(), expects: z.literal('option').default('option') })
        : poll.expects === 'multi'
          ? zPollAnswerSchema
              .omit({ answer: true })
              .extend({ answer: z.string().array(), expects: z.literal('multi').default('multi') })
          : zPollAnswerSchema
              .omit({ answer: true })
              .extend({ answer: z.number(), expects: z.literal('number').default('number') });
  const zErrorSchema = zMakeErrorDataSchema(zFormSchema.keyof());

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(
    defaults(zod(zFormSchema)),
    {
      SPA: true,
      validators: zod(zFormSchema),
      async onUpdate({ form }) {
        if (!form.valid) return;

        try {
          await pb.collection('pollAnswers').create({ ...form.data, poll: poll.id, user: user.id });
        } catch (e) {
          return processError(form, e, zErrorSchema);
        }
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
  forceDisable={!poll.active}
>
  <svelte:fragment slot="title">{poll.title}</svelte:fragment>

  {#if $form.expects === 'number'}
    <NumberInput
      name="answer"
      errors={$errors.answer}
      bind:value={$form.answer}
      constraints={$constraints.answer}
    >
      Resposta
    </NumberInput>
  {:else if $form.expects === 'option'}
    <SelectInput
      name="answer"
      errors={$errors.answer}
      bind:value={$form.answer}
      constraints={$constraints.answer}
      options={poll.options.split(',')}>Resposta</SelectInput
    >
  {:else if $form.expects === 'multi'}
    <div class="max-w-xs">
      <div class="label">Resposta</div>
      {#each poll.options.split(',') as opt}
        <div class="form-control">
          <label class="label cursor-pointer">
            <span class="label-text">{opt}</span>
            <input type="checkbox" value={opt} bind:group={$form.answer} class="checkbox" />
          </label>
        </div>
      {/each}
    </div>
  {:else if $form.expects === 'string'}
    <TextInput
      name="answer"
      errors={$errors.answer}
      bind:value={$form.answer}
      constraints={$constraints.answer}
    >
      Resposta
    </TextInput>
  {:else}
    <ErrorAlert errors={['Something is wrong']} />
  {/if}

  <svelte:fragment slot="save">Salvar</svelte:fragment>
</BasicFormCard>
