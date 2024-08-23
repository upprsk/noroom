<script>
  import BasicFormCard from '$lib/components/BasicFormCard.svelte';
  import { filesProxy, superForm } from 'sveltekit-superforms';

  export let data;

  const { form, errors, constraints, message, submitting, delayed, enhance } = superForm(data.form);
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

  <div class="flex w-full gap-2 flex-wrap">
    {#each data.klass.attachments as attach}
      <div class="badge badge-ghost">
        <svg class="w-4 h-4 fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256"
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
      bind:value={$files}
      aria-invalid={$errors.attachments ? 'true' : undefined}
      class="file-input file-input-bordered w-full max-w-xs"
      {...constraints}
    />
    <div class="label">
      {#if $errors.attachments}
        <span class="label-text-alt text-error">
          {Object.entries($errors.attachments).map(([k, v]) => `in ${k}: ${v}`)}
        </span>
      {/if}
    </div>
  </label>

  <svelte:fragment slot="actions">
    <button type="button" class="btn" on:click={() => history.back()}>voltar</button>
  </svelte:fragment>
</BasicFormCard>
