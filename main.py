import os
import time
import discord
from discord.ext import commands
from colorama import Fore
from dotenv import load_dotenv

load_dotenv()

PREFIX = os.getenv("PREFIX")
TOKEN = os.getenv("TOKEN")

if not PREFIX or not TOKEN:
    raise ValueError("PREFIX or TOKEN is not set")

bot = commands.Bot(command_prefix=PREFIX, help_command=None, self_bot=True)

@bot.event
async def on_ready():
    print(f"[Logged in as {bot.user}]\nLatency: {bot.latency*1000}ms")

@bot.command()
async def purge(ctx, channel_id: int, amount: int, limit: float):
    count = 0
    channel = bot.get_channel(int(channel_id))
    if not channel:
        await ctx.reply(f"Channel not found: `{channel_id}`")
        return
    messages = [message async for message in channel.history(limit=amount + 1)]
    for msg in messages:
        if msg.author == bot.user:
            try:
                await msg.delete()
                count += 1
                print(Fore.RED+"[DELETED]"+Fore.RESET+f" {msg.author} | {msg.content}")
                time.sleep(float(limit))
            except discord.errors.Forbidden as e:
                if e.code == 50021: 
                    amount += 1
                    pass
        else:
            amount += 1
    await ctx.send(f"`deleted`", delete_after=5)
    print(Fore.GREEN+"[FINISHED]"+Fore.RESET+f" Count: {count}")

@purge.error
async def purge_error(ctx, error):
    if isinstance(error, commands.MissingRequiredArgument):
        usage = f"```{PREFIX}purge [channel_id] [amount] [float(time)]\n{PREFIX}purge 1370064823085170698 100 1.45```"
        await ctx.reply(usage, delete_after=10)
    else:
        raise error

if __name__ == "__main__":
    bot.run(TOKEN)