<script lang="ts">
  import Navbar from '$lib/components/Navbar.svelte';
  import { currentUser } from '$lib/stores/user';
  import '../app.css';

  export let data;

  // Set the current user from the data passed in from the server
  $: currentUser.set(data.user);
</script>

<div class="flex flex-col w-full h-full">
  <Navbar />

  <div class="grow overflow-y-auto">
    <slot />
  </div>
</div>
