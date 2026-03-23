using StardewModdingAPI;
namespace InventoryManager;

public class ModEntry : Mod
{
  public override void Entry(IModHelper helper)
  {
    Monitor.Log("Mod loaded!", LogLevel.Info);
  }
}