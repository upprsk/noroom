<script lang="ts">
  import { page } from '$app/stores';
  import Alert from './Alert.svelte';

  export let enhance;
  export let delayed: boolean;
  export let submitting: boolean;
  export let _errors: string[] | undefined;
  export let message: string | undefined;
  export let forceDisable = false;
</script>

<div class="flex items-center justify-center p-5">
  <div class="card w-full max-w-4xl bg-base-100">
    <form class="card-body" method="POST" enctype="multipart/form-data" use:enhance>
      <h4 class="card-title"><slot name="title" /></h4>

      <slot />

      {#if _errors}
        {#each _errors as err}
          <Alert class="alert-error">
            <span>{err}</span>
          </Alert>
        {/each}
      {/if}

      {#if message}
        <Alert
          class="{$page.status === 200 ? 'alert-success' : ''} {$page.status >= 400
            ? 'alert-error'
            : ''}"
        >
          <span>{message}</span>
        </Alert>
      {/if}

      <slot name="notes" />

      <div class="card-actions justify-end">
        <slot name="actions" />

        <button class="btn btn-primary" disabled={submitting || forceDisable}>
          {#if delayed}
            <span class="loading loading-spinner loading-sm"></span>
          {:else}
            <slot name="save">
              <slot name="title" />
            </slot>
          {/if}
        </button>
      </div>
    </form>
  </div>
</div>
