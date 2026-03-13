from typing import TypedDict, NotRequired
from probability import distribute_balls_unlabeled, unique_permutation_count
import xml.etree.ElementTree as ET

class Item(TypedDict):
  name:str
  quantity:int

# encodes an item into a string for easy comparisons
def encode_item(item:Item):
  return f"{item['name']}//{item['quantity']}"

def decode_item(string:str)->Item:
  parts = string.split("//")
  item:Item =  {
    "name": parts[0],
    "quantity": int(parts[1]),
  }
  return item

# given a full encoded inventory, figure out the different arrangements and return them as encoded items
def calculate_inventory_arrangements(inventory:list[str]):
  return unique_permutation_count(inventory)

# given a stack size, and number of slots to spread it to, return all the possible combinations (order ignored)
def spread_stack(stack:int, slots:int):
  return distribute_balls_unlabeled(stack, slots)

# TODO handle multiple stacks of different items
def calculate_inventories(items:list[Item]):
  de_blanked = [item for item in items if item["name"] != "Blank"]
  item_names = [item["name"] for item in de_blanked]
  if len(de_blanked) != len(set(item_names)):
    raise Exception("No duplicate items allowed")
  inventories = {",".join([encode_item(item) for item in items])}
  
  blanks = len(items) - len(de_blanked)
  for item in de_blanked:
    if(item["quantity"]<2):
      continue
    for i in range(blanks):
      slots = i+2
      stack = item["quantity"]
      spreads = spread_stack(stack, slots)
      for spread in spreads:
        inventories.add(",".join(sorted(
          [encode_item(spread_item) for spread_item in de_blanked if (spread_item["name"] != item["name"])] +
          [f"{item['name']}//{spread_count}//False" for spread_count in spread]+
          ["Blank//1//True" for _ in range(blanks-i-1)]
        )))
  return inventories

def get_inventory_counts(items:list[Item], verbose=False):
  # The different combinations of items
  encoded_inventories = sorted(calculate_inventories(items))
  inventory_counts = {}
  abstract_inventories = {}

  # For each inventory, calculate the number of permutations for it
  for i,encoded_inventory in enumerate(encoded_inventories):
    if verbose:
      print(f"Checking Inventory {i+1}/{len(encoded_inventories)}")
    inventory = encoded_inventory.split(",")
    abstract_inventory_dict = {}
    for item in inventory:
      if item not in abstract_inventory_dict.keys():
        abstract_inventory_dict[item] = 0
      abstract_inventory_dict[item] += 1
    abstract_inventory = ",".join(str(v) for v in [sorted(abstract_inventory_dict.values())])
    if abstract_inventory in abstract_inventories.keys():
      i_count = abstract_inventories[abstract_inventory]
    else:
      i_count = calculate_inventory_arrangements(inventory)
      abstract_inventories[abstract_inventory] = i_count
    if verbose:
      print(f"Found {i_count} permutations")
    inventory_counts[encoded_inventory] = i_count
  return inventory_counts


def get_inventory(number:int, items:list[Item], verbose=False) -> list[Item]:
  # Get all the inventory counts
  inventory_counts = get_inventory_counts(items)
  count = 0
  encoded_inventory = None
  # For each inventory, see how many permutations it has until you have overshot the target inventory
  for key in inventory_counts.keys():
    cur_count = inventory_counts[key]
    if count+cur_count>=number:
      encoded_inventory = str(key)
      break
    count += cur_count
  if not encoded_inventory:
    if verbose:
      print("No valid inventory found")
    return []
  
  inventory_combination = encoded_inventory.split(",")
  if verbose:
    print(f"Found inventory combination: {inventory_combination}")
  final_inventory = []
  perms_to_check = [i for i in inventory_combination]
  while len(perms_to_check) > 0:
    checked_items = set()
    for item in perms_to_check:
      if item in checked_items:
        continue
      checked_items.add(item)
      remaining_inventory = [i for i in perms_to_check]
      remaining_inventory.remove(item)
      perm_count = unique_permutation_count(remaining_inventory)
      if count+perm_count >= number:
        final_inventory.append(decode_item(item))
        perms_to_check.remove(item)
        break
      count += perm_count
  return final_inventory

def get_inventory_index(inventory:list[Item], items:list[Item]):
  sorted_inventory = sorted([encode_item(i) for i in inventory])
  encoded_sorted_inventory = ",".join(sorted_inventory)
  inventory_counts = get_inventory_counts(items)
  if encoded_sorted_inventory not in inventory_counts.keys():
    return -1
  count = 1
  for key in inventory_counts:
    if key == encoded_sorted_inventory:
      break
    count += inventory_counts[key]

  perms_to_check = [i for i in sorted_inventory]
  for target_item in [encode_item(i) for i in inventory]:
    checked_items = set()
    for item in perms_to_check:
      if item in checked_items:
        continue
      checked_items.add(item)
      remaining_inventory = [i for i in perms_to_check]
      remaining_inventory.remove(item)
      perm_count = unique_permutation_count(remaining_inventory)
      if item == target_item:
        perms_to_check.remove(item)
        break
      count += perm_count
  return count
  


def text_to_number(message:str):
  binary_message = [f"{ord(c):0{8}b}" for c in message]
  if len(binary_message)>32:
    return -1
  return int("".join(binary_message),2)


def number_to_text(number:int):
  binary_string = f"{number:032b}"
  return "".join(chr(int(binary_string[i:i+8], 2)) for i in range(0, len(binary_string), 8))

def get_current_inventory(save_path:str):
  tree = ET.parse(save_path)
  root = tree.getroot()
  items = root.find("items")
  inventory:list[Item] = []
  if items is None:
    return []
  
  for item in items:
    name = item.find("name")
    quantity = item.find("stack")
    if name is None or quantity is None:
      inventory.append({
        "name":"Blank",
        "quantity":1
      })
    else:
      inventory.append({
        "name":str(name.text),
        "quantity":int(str(quantity.text))
      })

  return inventory
  
  

'''
# 5157250560
print(f"Total: {sum(get_inventory_counts(items).values())}")
print(f"Accuracy: {sum(get_inventory_counts(items).values())/51572505.6}%")
# print(get_inventory(text_to_number("Hi!"),items))

# Encode
msg = "hey!"
number_msg = text_to_number(msg)
inventory = get_inventory(number_msg,items)

inventory_index = get_inventory_index(inventory,items)
decoded_msg = number_to_text(inventory_index)

print(f"Encrypting Message: {msg}")
print(f"Convert to number: {number_msg}")
print(f"Found inventory:\n - {'\n - '.join([str(item['quantity'])+item['name'] for item in inventory])}")

print(f"Inventory Index: {inventory_index}")
print(f"Decoded message: {decoded_msg}")
'''

print(get_current_inventory("/Users/ocapraro/.config/StardewValley/Saves/CHANNEL_431325361/SaveGameInfo"))