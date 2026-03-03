from typing import TypedDict, NotRequired
from itertools import permutations
from probability import distribute_balls_unlabeled, unique_permutation_count

class Item(TypedDict):
  name:str
  quantity:int
  unstackable:NotRequired[bool]

# encodes an item into a string for easy comparisons
def encode_item(item:Item):
  return f"{item['name']}//{item['quantity']}//{item.get("unstackable", False)}"

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
        inventories.add(",".join(
          [encode_item(spread_item) for spread_item in de_blanked if (spread_item["name"] != item["name"])] +
          [f"{item['name']}//{spread_count}" for spread_count in spread]+
          ["Blank//1//True" for _ in range(blanks-i-1)]
        ))
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


def get_inventory(number:int, items:list[Item], verbose=False):
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
  perm_number = number - count
  final_inventory = []
  perms_to_check = [i for i in inventory_combination]
  count = 0
  while len(perms_to_check) > 0:
    checked_items = set()
    for item in perms_to_check:
      if item in checked_items:
        continue
      checked_items.add(item)
      remaining_inventory = [i for i in perms_to_check]
      remaining_inventory.remove(item)
      perm_count = unique_permutation_count(remaining_inventory)
      if count+perm_count >= perm_number:
        final_inventory.append(item)
        perms_to_check.remove(item)
        break
      count += perm_count
  return final_inventory


items:list[Item] = [
  {
    "name":"Axe",
    "quantity":1
  },
  {
    "name":"Hoe",
    "quantity":1
  },
  {
    "name":"WateringCan",
    "quantity":1
  },
  {
    "name":"Pickaxe",
    "quantity":1
  },
  {
    "name":"MeleeWeapon",
    "quantity":1
  },
  {
    "name":"Parsnip Seeds",
    "quantity":15
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  },
  {
    "name":"Blank",
    "quantity":1,
    "unstackable":True
  }
]



print(f"Total: {sum(get_inventory_counts(items).values())}")
print("\n".join(get_inventory(10000, items)))