using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Xml;

namespace WxiFindId
{
    class Program
    {
        static void Main(string[] args)
        {
            if (args.Length !=  3)
            {
                DisplayUsage();
            }

            XmlDocument xml = new XmlDocument();
            try
            {
                xml.Load(args[2]);
            }
            catch (XmlException ex)
            {
                Console.WriteLine(ex.Message);
                //Cannot open xml file
                System.Environment.Exit(2);
            }

            string nodeType;
            string attributeName;
            if (args[0] == "-d")
            {
                nodeType = "Directory";
                attributeName = "Name";
            }
            else if (args[0] == "-f")
            {
                nodeType = "File";
                attributeName = "Source";
            }
            else if (args[0] == "-c")
            {
                nodeType = "Component";
                attributeName = "Id";
            }
            else
            {
                DisplayUsage();
                return;
            }

            XmlNodeList xnList = xml.GetElementsByTagName(nodeType);
            foreach (XmlNode xn in xnList)
            {
                try
                {
                    XmlElement xe = (XmlElement)xn;
                    if (xe.HasAttribute(attributeName))
                    {
                        XmlAttribute xa = xe.GetAttributeNode(attributeName);
                        if (xa.Value.EndsWith(args[1]))
                        {
                            if (xe.HasAttribute("Id"))
                            {
                                XmlAttribute xaId = xe.GetAttributeNode("Id");
                                //First occurence found, write to stdout and exit 0
                                Console.Write(xaId.Value);
                                System.Environment.Exit(0);
                            }
                        }
                    }
                }
                catch { }
            }
            //Not found
            System.Environment.Exit(3);
        }

        static void DisplayUsage()
        {
            Console.WriteLine("Usage: WxiFindId [-d|f|c] searchName FilesList.wxi");
            Console.WriteLine(" -d return Id a first Directory node with Name ending with searchName");
            Console.WriteLine(" -f return Id of first File node with Source ending with searchName");
            Console.WriteLine(" -c return Id of the component associated with first File node which Source ending with searchName");
            //Wrong arguments
            System.Environment.Exit(1);
        }
    }
}
