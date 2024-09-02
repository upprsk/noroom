<script>
  import { page } from '$app/stores';
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import { zErrorSchema } from '$lib/models';
  import { pb, processError } from '$lib/pocketbase';
  import { defaults, filesProxy, setMessage, superForm } from 'sveltekit-superforms';
  import { zod } from 'sveltekit-superforms/adapters';
  import { zFormSchema } from './models';

  export let data;

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(
    defaults(zod(zFormSchema)),
    {
      SPA: true,
      validators: zod(zFormSchema),
      async onUpdate({ form }) {
        if (!form.valid) return;

        try {
          await pb.collection('classes').update($page.params.id, form.data);
        } catch (e) {
          return processError(form, e, zErrorSchema);
        }

        return setMessage(form, `${form.data.attachments.length} arquivos anexados`);
      },
    },
  );

  const files = filesProxy(form, 'attachments');
</script>

<BasicFormCard
  {enhance}
  delayed={$delayed}
  submitting={$submitting}
  _errors={$errors._errors}
  message={$message}
>
  <svelte:fragment slot="title">Anexar arquivos</svelte:fragment>

  <div class="flex w-full flex-wrap gap-2">
    {#each data.klass.attachments as attach}
      <div class="badge badge-ghost">
        <svg class="h-4 w-4 fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256"
          ><path
            d="M213.66,82.34l-56-56A8,8,0,0,0,152,24H56A16,16,0,0,0,40,40V216a16,16,0,0,0,16,16H200a16,16,0,0,0,16-16V88A8,8,0,0,0,213.66,82.34ZM160,51.31,188.69,80H160ZM200,216H56V40h88V88a8,8,0,0,0,8,8h48V216Z"
          ></path></svg
        >
        {attach}
      </div>
    {/each}
  </div>

  <label class="form-control w-full max-w-xs">
    <div class="label">
      <span class="label-text"><slot /></span>
    </div>
    <input
      multiple
      type="file"
      name="attachments"
      bind:files={$files}
      aria-invalid={$errors.attachments ? 'true' : undefined}
      class="file-input file-input-bordered w-full max-w-xs"
      {...constraints}
    />
    <div class="label">
      {#if $errors.attachments}
        <span class="label-text-alt text-error">
          {#if $errors.attachments._errors}
            {$errors.attachments._errors}
          {/if}

          <!-- {Object.entries($errors.attachments).map(([k, v]) => `in ${k}: ${v}`)} -->
        </span>
      {/if}
    </div>
  </label>

  <svelte:fragment slot="actions">
    <button type="button" class="btn" on:click={() => history.back()}>voltar</button>
  </svelte:fragment>
</BasicFormCard>
