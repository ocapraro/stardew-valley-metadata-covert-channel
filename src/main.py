from typing import TypedDict, NotRequired
from itertools import permutations
from probability import distribute_balls_unlabeled

class FileReaderPOC:
  def loadSave(self):
    pass

  def getSteps(self):
    pass

class Item(TypedDict):
  name:str
  quantity:int
  unstackable:NotRequired[bool]


def are_inventories_equal(i1:list[Item], i2:list[Item]):
  if len(i1) != len(i2):
    return False
  for i, item1 in enumerate(i1):
    item2 = i2[i]
    if item1["name"] != item2["name"]:
      return False
    if item1["quantity"] != item2["quantity"]:
      return False

# encodes an item into a string for easy comparisons
def encode_item(item:Item):
  return f"{item['name']}//{item['quantity']}//{item.get("unstackable", False)}"

# given a full encoded inventory, figure out the different arrangements and return them as encoded items
def calculate_inventory_arrangements(inventory:list[str]):
  return set(permutations(inventory))

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
      
  
  # foreach item look to see if it appears in the rest of the list
  # if it does 
  # remove a blank



# Given a list of items, each with a quanity that can ba split into different stacks, generate the total number of inventories possible with a 10 slot inventory
# def generateInventoryCombinations(items:list[Item]) -> int:
#   max_slots = 10
  # calculate total number of rearrangements
  # calculate total number of stack

  

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
  }
]

count = 0
inventories = calculate_inventories(items)
# print(inventories.pop().split(","))
for ii,i in enumerate(inventories):
  print(f"Checking Inventory {ii+1}/{len(inventories)}")
  inventory = i.split(",")
  i_count = len(calculate_inventory_arrangements(inventory))
  print(f"Found {i_count} permutations")
  count += i_count
print(f"Total: {count}")