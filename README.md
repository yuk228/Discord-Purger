# Discord-Purger

## Overview

This tool is an automation script for bulk deleting your own messages in a specified Discord channel.  
It is created for educational purposes only and use is at your own risk.

## Setup Instructions

1. Install required libraries  
   ```bash
   pip install -r requirements.txt
   ```
2. Rename `.env.example` to `.env` and set environment variables
3. Run the script  
   ```bash
   python main.py
   ```

## Usage

Command format:
```
-purge [channel_id] [amount] [float(time)]
```

- `channel_id` : ID of the channel where you want to delete messages
- `amount` : Number of messages to delete (integer)
- `float(time)` : Deletion interval (seconds, decimal. Recommended: 1.45)

**Example:**
```
-purge 1370064823085170698 100 1.45
```

## Disclaimer

This program is created for **educational purposes only**. **Use is completely at your own risk**, and the developer assumes no responsibility for any damages or issues that may arise from the use of this program. Discord account automation may violate Discord's Terms of Service, and I do not recommend execution. Please use only within the scope of hobbies.
