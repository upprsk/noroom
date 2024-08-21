<script lang="ts">
  import BasicCard from '$lib/components/BasicCard.svelte';
  import EmailInput from '$lib/components/input/EmailInput.svelte';
  import IntInput from '$lib/components/input/IntInput.svelte';
  import TextInput from '$lib/components/input/TextInput.svelte';
  import { zEndDeviceSchema } from '$lib/models.js';
  import { pb, updateFromEvent } from '$lib/pocketbase.js';
  import { onMount } from 'svelte';
  import type { z } from 'zod';

  type EndDeviceModel = z.infer<typeof zEndDeviceSchema>;

  export let data;

  let unsub: () => void;
  const sub = async () => {
    unsub = await pb
      .collection('endDevices')
      .subscribe<EndDeviceModel>(
        '*',
        (e) => (data.devices = updateFromEvent(e, zEndDeviceSchema, data.devices)),
      );
  };

  onMount(() => {
    sub();

    return () => {
      if (unsub) unsub();
    };
  });
</script>

<BasicCard>
  <svelte:fragment slot="title">
    {#if data.user}
      {data.user.name} - {data.user.mat} - {data.user.curso}
    {:else}
      ...
    {/if}
  </svelte:fragment>
  {#if data.user}
    <h5 class="text-slate-400 text-xs italic">@{data.user.username} - {data.user.email}</h5>

    <div class="h-2"></div>

    <p>
      Ao trocar email ou senha, um email automatico vai ser enviado para {data.user.email} com instrucoes.
    </p>

    <div class="flex gap-2">
      <button class="btn btn-sm btn-warning">trocar email</button>
      <button class="btn btn-sm btn-warning">trocar senha</button>
    </div>

    <div class="divider"></div>

    <h4 class="card-title">Seus dispositivos</h4>
    <ul class="mx-5">
      {#each data.devices as dev (dev.id)}
        <li class="border-b last:border-b-0 py-5">
          <div>
            <span>
              <b>{dev.deviceData?.system.browser.name} {dev.deviceData?.system.browser.version}</b>
              -
              <b>{dev.deviceData?.system.platform}</b>
            </span>
            <span class="text-sm">{dev.locationData?.city}</span>
          </div>
        </li>
      {/each}
    </ul>
  {:else}
    ...
  {/if}
</BasicCard>
