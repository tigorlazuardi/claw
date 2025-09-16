<script lang="ts">
  import { useQueryClient } from "@sveltestack/svelte-query";
  import { getSourceServiceClient } from "../../connectrpc";
  import { toDate } from "../../connectrpc/js_date";
  import { parseCronExpression } from "cron-schedule";
  import { Popover } from "bits-ui";
  import IconInfo from "@lucide/svelte/icons/info";

  const queryClient = useQueryClient();

  interface Props {
    value: string[];
  }

  let { value = $bindable([]) }: Props = $props();

  async function getNextRun(pattern: string) {
    return queryClient.fetchQuery(
      ["schedules", "nextRun", pattern],
      () =>
        getSourceServiceClient().then((client) =>
          client.getCronNextTime({
            cronExpression: pattern,
          }),
        ),
      {
        staleTime: 60 * 1000,
        cacheTime: 60 * 1000,
        retry: false,
      },
    );
  }

  type ScheduleEntry = {
    pattern: string;
    nextRun?: Date;
    error?: string;
    zone?: string;
  };

  const schedules = $derived.by(async () => {
    const result: ScheduleEntry[] = [];
    for (const pattern of value) {
      const entry: ScheduleEntry = { pattern };
      try {
        parseCronExpression(pattern);
        const resp = await getNextRun(pattern);
        entry.nextRun = toDate(resp.nextTime!);
      } catch (err) {
        if (err instanceof Error) {
          entry.error = err.message;
        } else {
          entry.error = "Invalid cron expression";
        }
      }
      result.push(entry);
    }
    return result;
  });

  type PatternOption = {
    label: string;
    pattern: string;
    description: string;
  };

  type ScheduleInputState = {
    value: string;
    element?: HTMLInputElement;
    nextRun?: Date;
    error?: string;
    zone?: string;
  };

  let scheduleInputState: ScheduleInputState = $state({ value: "" });
  const cleanedValue = $derived(scheduleInputState.value.trim());

  async function handleOnInput() {
    scheduleInputState.element?.setCustomValidity("");
    scheduleInputState.nextRun = undefined;
    scheduleInputState.error = undefined;
    scheduleInputState.zone = undefined;
    if (!cleanedValue) {
      return;
    }
    try {
      parseCronExpression(cleanedValue);
      queryClient.invalidateQueries(["schedules", "nextRun", cleanedValue]);
      const resp = await getNextRun(cleanedValue);
      scheduleInputState.nextRun = toDate(resp.nextTime!);
      scheduleInputState.zone = resp.zone || "UTC";
    } catch (err) {
      if (err instanceof Error) {
        scheduleInputState.element?.setCustomValidity(err.message);
        scheduleInputState.error = err.message;
      } else {
        scheduleInputState.element?.setCustomValidity(`${err}`);
        scheduleInputState.error = `Invalid cron expression: ${err}`;
      }
    }
  }

  const patternDropdownList: PatternOption[] = [
    {
      label: "Every minute",
      pattern: "* * * * *",
      description: "Runs every minute",
    },
    {
      label: "Every minute (alt)",
      pattern: "@minutely",
      description: "Runs every minute",
    },
    {
      label: "Every 5 minutes",
      pattern: "*/5 * * * *",
      description: "Runs every 5 minutes",
    },
    {
      label: "Every 15 minutes",
      pattern: "*/15 * * * *",
      description: "Runs every 15 minutes",
    },
    {
      label: "Every 30 minutes",
      pattern: "*/30 * * * *",
      description: "Runs every 30 minutes",
    },
    {
      label: "Every hour",
      pattern: "0 * * * *",
      description: "Runs at the start of every hour",
    },
    {
      label: "Daily at midnight",
      pattern: "0 0 * * *",
      description: "Runs daily at midnight",
    },
    {
      label: "Weekly on Sunday at midnight",
      pattern: "0 0 * * 0",
      description: "Runs weekly on Sunday at midnight",
    },
    {
      label: "Monthly on the 1st at midnight",
      pattern: "0 0 1 * *",
      description: "Runs monthly on the 1st at midnight",
    },
  ];
</script>

<fieldset class="fieldset">
  <legend class="fieldset-legend">
    <span>Schedule</span>
  </legend>
  <input
    name="schedule"
    bind:this={scheduleInputState.element}
    type="text"
    class={{
      "input w-full": true,
    }}
    bind:value={scheduleInputState.value}
    placeholder="e.g. 0 0 * * FRI (every midnight at Friday)"
    oninput={handleOnInput}
  />
  {#if scheduleInputState.nextRun}
    {@const locale = Intl.NumberFormat().resolvedOptions().locale}
    <div class="alert alert-success alert-soft mt-2">
      <span>
        Next run time: {scheduleInputState.nextRun.toLocaleString(locale)} (adjusted
        to current timezone {Intl.DateTimeFormat().resolvedOptions().timeZone}).
        Server time: {scheduleInputState.nextRun.toLocaleString(locale, {
          timeZone: scheduleInputState.zone,
        })} ({scheduleInputState.zone}).
      </span>
    </div>
  {:else if queryClient.isFetching(["schedules", "nextRun"])}
    <div class="alert alert-error alert-soft">
      <span class="loading loading-spinner"></span>
      <span>Checking next run time on the server...</span>
    </div>
  {:else if scheduleInputState.error}
    <div class="alert alert-error alert-soft">
      <span>{scheduleInputState.error}</span>
    </div>
  {:else}
    <p class="label">Schedule using Cron Expression</p>
  {/if}
</fieldset>

{#snippet helpText()}
  <Popover.Root>
    <Popover.Trigger class="btn btn-square btn-ghost btn-xs">
      <IconInfo />
    </Popover.Trigger>
    <Popover.Portal>
      <Popover.Content class="z-[9999] bg-transparent" data-theme="dracula"
      ></Popover.Content>
    </Popover.Portal>
  </Popover.Root>
{/snippet}
