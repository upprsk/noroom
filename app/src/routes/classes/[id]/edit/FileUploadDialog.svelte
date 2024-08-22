<script lang="ts">
  import { page } from '$app/stores';
  import Alert from '$lib/components/Alert.svelte';
  import { zFileUploadSchema } from '$lib/models';
  import { defaults, superForm } from 'sveltekit-superforms';
  import { zod } from 'sveltekit-superforms/adapters';

  export const show = () => {
    dialog.show();
  };

  export const close = () => {
    dialog.close();
  };

  const { form, message, enhance, errors } = superForm(defaults(zod(zFileUploadSchema)), {
    onUpdated({ form }) {
      if (form.valid) {
        console.log(form.data);
      }
    },
  });

  let dialog: HTMLDialogElement;
</script>

<dialog bind:this={dialog} class="modal modal-bottom sm:modal-middle">
  <form class="modal-box" method="POST" action="?/upload" enctype="multipart/form-data" use:enhance>
    <h3 class="text-lg font-bold">Anexar arquivos</h3>

    <label class="form-control w-full max-w-xs">
      <div class="label">
        <span class="label-text">Anexar arquivos</span>
      </div>
      <input
        multiple
        type="file"
        name="attachments"
        class="file-input file-input-bordered w-full max-w-xs"
        on:input={(e) => {
          if (!e.currentTarget.files) {
            $form.files = null;
            return;
          }

          const files = [];
          for (const file of e.currentTarget.files) {
            files.push(file);
          }

          $form.files = files;
        }}
      />
      <div class="label">
        {#if $errors.files}
          <span class="label-text-alt text-error">
            {Object.entries($errors.files).map(([k, v]) => `at ${k}: ${v}`)}
          </span>
        {/if}
      </div>
    </label>

    {#if $errors._errors}
      {#each $errors._errors as err}
        <Alert class="alert-error">
          <span>{err}</span>
        </Alert>
      {/each}
    {/if}

    {#if $message}
      <Alert
        class="{$page.status === 200 ? 'alert-success' : ''} {$page.status >= 400
          ? 'alert-error'
          : ''}"
      >
        <span>{$message}</span>
      </Alert>
    {/if}

    <div class="modal-action">
      <button type="submit" class="btn">Anexar</button>
    </div>
  </form>
</dialog>
