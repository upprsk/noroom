<script lang="ts">
  import { currentUser } from '$lib/stores/user';
  import BasicAvatar from './BasicAvatar.svelte';
</script>

<div class="navbar-end">
  {#if $currentUser}
    <a href="/account">
      <BasicAvatar user={$currentUser} />
    </a>
  {:else}
    <a href="/login" class="btn-sm btn btn-accent">Sign-In</a>
  {/if}
</div>
