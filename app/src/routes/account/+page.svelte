<script lang="ts">
  import BasicCard from '$lib/components/BasicCard.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import EmailInput from '$lib/components/input/EmailInput.svelte';
  import StatusAlert from '$lib/components/StatusAlert.svelte';
  import { zEndDeviceSchema } from '$lib/models.js';
  import { pb, updateFromEvent } from '$lib/pocketbase.js';
  import { ClientResponseError } from 'pocketbase';
  import { onMount } from 'svelte';
  import type { z } from 'zod';

  type EndDeviceModel = z.infer<typeof zEndDeviceSchema>;

  export let data;

  let changeEmailErrors: undefined | string[];
  let changeEmailStatus: undefined | string;
  let changeEmailValue = '';
  const changeEmail = async () => {
    changeEmailErrors = undefined;
    changeEmailStatus = undefined;

    try {
      await pb.collection('users').requestEmailChange(changeEmailValue);

      changeEmailStatus = 'Modificacao requisitada. Verifique sua caixa de entrada';
    } catch (e) {
      if (e instanceof ClientResponseError) {
        changeEmailErrors = [`Erro ao requisitar troca: ${e.message}`];
        console.log(e.data);
      } else {
        changeEmailErrors = [`Erro ao requisitar troca: ${e}`];
      }
    }
  };

  // --------------------------------------------------------------------------

  let confirmAccountErrors: undefined | string[];
  let confirmAccountStatus: undefined | string;
  const confirmAccount = async () => {
    confirmAccountErrors = undefined;
    confirmAccountStatus = undefined;

    try {
      await pb.collection('users').requestVerification(data.user.email!);

      confirmAccountStatus = 'Confirmacao requisitada. Verifique sua caixa de entrada';
    } catch (e) {
      if (e instanceof ClientResponseError) {
        confirmAccountErrors = [`Erro ao requisitar troca: ${e.message}`];
        console.log(e.data);
      } else {
        confirmAccountErrors = [`Erro ao requisitar troca: ${e}`];
      }
    }
  };

  // --------------------------------------------------------------------------

  let uploadAvatarErrors: undefined | string[];
  let uploadAvatarErrorsExtra: Record<string, string> = {};
  let uploadAvatarStatus: undefined | string;
  let file: FileList | undefined;
  const uploadAvatar = async () => {
    uploadAvatarErrors = undefined;
    uploadAvatarStatus = undefined;
    uploadAvatarErrorsExtra = {};

    try {
      const res = await pb.collection('users').update(data.user.id, {
        avatar: file?.item(0) ?? null,
      });

      if (res.avatar !== '') uploadAvatarStatus = 'Avatar alterado';
      else uploadAvatarStatus = 'Avatar removido';
    } catch (e) {
      if (e instanceof ClientResponseError) {
        uploadAvatarErrors = [`Erro ao atualizar avatar: ${e.message}`];
        console.log(e.data);

        for (const [k, v] of Object.entries(e.data.data)) {
          const { message } = v as { message: string };
          uploadAvatarErrorsExtra[k] = message;
        }
      } else {
        uploadAvatarErrors = [`Erro ao atualizar avatar: ${e}`];
      }
    }
  };

  // --------------------------------------------------------------------------

  let resetPasswordErrors: undefined | string[];
  const resetPassword = async () => {
    resetPasswordErrors = ['nao implementado'];
  };

  // --------------------------------------------------------------------------

  let unsub: () => void;
  const sub = async () => {
    unsub = await pb
      .collection('endDevices')
      .subscribe<EndDeviceModel>(
        '*',
        (e) => (data.devices = updateFromEvent(e, zEndDeviceSchema, data.devices)),
      );
  };

  // --------------------------------------------------------------------------

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
    <h5 class="text-xs italic text-slate-400">@{data.user.username} - {data.user.email}</h5>

    <div class="h-2"></div>

    {#if !data.user.verified}
      <div role="alert" class="alert alert-warning">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6 shrink-0 stroke-current"
          fill="none"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
        <span>Sua conta ainda nao foi validada!</span>
        <button type="button" class="btn btn-sm" on:click={confirmAccount}>validar</button>
      </div>

      <div class="h-2"></div>

      <StatusAlert status={confirmAccountStatus} />
      <ErrorAlert errors={confirmAccountErrors} />

      <div class="h-2"></div>
    {/if}

    <div>
      <label class="form-control w-full max-w-xs">
        <div class="label">
          <span class="label-text">Avatar</span>
        </div>
        <input
          name="avatar"
          type="file"
          class="file-input file-input-bordered w-full max-w-xs"
          bind:files={file}
        />
        <div class="label">
          {#if uploadAvatarErrorsExtra.avatar}
            <span class="label-text-alt text-error">{uploadAvatarErrorsExtra.avatar}</span>
          {/if}
        </div>
      </label>

      <div class="h-2"></div>

      <StatusAlert status={uploadAvatarStatus} />
      <ErrorAlert errors={uploadAvatarErrors} />

      <div class="h-2"></div>

      <button type="button" class="btn" on:click={uploadAvatar}>
        {#if file?.length}atualizar{:else}remover{/if} avatar
      </button>
    </div>

    <div class="divider"></div>

    <div>
      <p>
        Ao requisitar a troca de senha, um email automatico vai ser enviado para <b
          >{data.user.email}</b
        > com instrucoes.
      </p>

      <button type="button" class="btn btn-warning btn-sm" on:click={resetPassword}>
        trocar senha
      </button>

      <div class="h-2"></div>

      <ErrorAlert errors={resetPasswordErrors} />
    </div>

    <div class="divider"></div>

    <div>
      <EmailInput
        name="email"
        errors={changeEmailErrors}
        bind:value={changeEmailValue}
        constraints={undefined}
      >
        Email
      </EmailInput>

      <button
        type="button"
        class="btn btn-warning btn-sm"
        on:click={changeEmail}
        disabled={changeEmailValue.length === 0}
      >
        trocar email
      </button>

      <div class="h-2"></div>

      <StatusAlert status={changeEmailStatus} />
      <ErrorAlert errors={changeEmailErrors} />
    </div>

    <div class="divider"></div>

    <h4 class="card-title">Seus dispositivos</h4>
    <ul class="mx-5">
      {#each data.devices as dev (dev.id)}
        <li class="border-b py-5 last:border-b-0">
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
